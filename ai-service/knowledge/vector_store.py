"""Milvus向量存储 - 知识库检索后端（懒加载）"""
from config import get_settings
# pymilvus 在首次使用时惰性导入，避免未安装时阻塞启动

_store_instances = {}


class VectorStore:
    """Milvus向量数据库封装"""

    def __init__(self, collection_name: str = None):
        settings = get_settings()
        self.collection_name = collection_name or settings.milvus_collection_name
        self.dim = settings.embedding_dim
        self._connected = False

    def _connect(self):
        """建立连接（懒连接）"""
        if self._connected:
            return
        from pymilvus import connections
        settings = get_settings()
        connections.connect(
            alias="default",
            host=settings.milvus_host,
            port=settings.milvus_port,
        )
        self._ensure_collection()
        self._connected = True

    def _ensure_collection(self):
        """确保Collection存在，不存在则创建"""
        from pymilvus import Collection, FieldSchema, CollectionSchema, DataType, utility

        if utility.has_collection(self.collection_name):
            self.collection = Collection(self.collection_name)
            return

        fields = [
            FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
            FieldSchema(name="content", dtype=DataType.VARCHAR, max_length=65535),
            FieldSchema(name="chunk_index", dtype=DataType.INT64),
            FieldSchema(name="metadata", dtype=DataType.JSON),
            FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=self.dim),
        ]
        schema = CollectionSchema(fields, description="儿童健康知识库")
        self.collection = Collection(self.collection_name, schema)

        index_params = {
            "metric_type": "IP",
            "index_type": "IVF_FLAT",
            "params": {"nlist": 1024},
        }
        self.collection.create_index("embedding", index_params)
        self.collection.load()

    def insert(self, chunks: list, embeddings: list) -> list:
        """批量插入chunk和对应向量，返回主键ID列表"""
        self._connect()
        data = [
            [c["content"] for c in chunks],
            [c.get("chunk_index", 0) for c in chunks],
            [c.get("metadata", {}) for c in chunks],
            embeddings,
        ]
        result = self.collection.insert(data)
        return result.primary_keys

    def search(
        self,
        query_vector: list,
        top_k: int = 10,
        filter_expr: str = None,
    ) -> list:
        """
        向量相似度检索。

        Args:
            query_vector: 查询向量
            top_k: 返回数量
            filter_expr: Milvus过滤表达式，如 'metadata["age_range"] == "0-1岁"'

        Returns:
            [{"id": ..., "content": "...", "score": 0.95, "metadata": {...}}, ...]
        """
        self._connect()
        self.collection.load()

        search_params = {"metric_type": "IP", "params": {"nprobe": 16}}
        results = self.collection.search(
            data=[query_vector],
            anns_field="embedding",
            param=search_params,
            limit=top_k,
            expr=filter_expr,
            output_fields=["content", "chunk_index", "metadata"],
        )

        hits = []
        for hit in results[0]:
            hits.append({
                "id": hit.id,
                "content": hit.entity.get("content"),
                "chunk_index": hit.entity.get("chunk_index"),
                "metadata": hit.entity.get("metadata"),
                "score": hit.score,
            })
        return hits

    def count(self) -> int:
        """返回存储的向量总数"""
        self._connect()
        return self.collection.num_entities

    def drop(self):
        """删除Collection（测试用）"""
        from pymilvus import utility
        if utility.has_collection(self.collection_name):
            utility.drop_collection(self.collection_name)

    def flush(self):
        """确保数据落盘"""
        self._connect()
        self.collection.flush()


def get_vector_store(collection_name: str = None) -> VectorStore:
    """获取VectorStore实例（按collection_name缓存）"""
    global _store_instances
    settings = get_settings()
    name = collection_name or settings.milvus_collection_name
    if name not in _store_instances:
        _store_instances[name] = VectorStore(collection_name=name)
    return _store_instances[name]
