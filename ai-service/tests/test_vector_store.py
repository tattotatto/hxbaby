"""Milvus向量存储测试
注意：需要Docker中Milvus运行。无Milvus时自动跳过。
"""
import pytest

COLLECTION = "test_knowledge_chunks"

# 检查Milvus是否可用
try:
    from pymilvus import connections, utility
    connections.connect(host="localhost", port=19530, timeout=3)
    _milvus_available = True
    connections.disconnect("default")
except Exception:
    _milvus_available = False

pytestmark = pytest.mark.skipif(
    not _milvus_available,
    reason="Milvus not running. Start with: docker compose up -d"
)


class TestVectorStore:
    @pytest.fixture(autouse=True)
    def setup(self):
        from knowledge.vector_store import get_vector_store
        self.store = get_vector_store(collection_name=COLLECTION)
        self.store.drop()
        # 重新创建
        self.store._connected = False
        yield
        self.store.drop()

    def test_insert_and_search(self):
        """插入并检索"""
        chunks = [
            {
                "content": "婴儿发热38.5度以下可物理降温",
                "chunk_index": 0,
                "metadata": {"category": "发热", "age_range": "0-1岁"},
            },
            {
                "content": "儿童腹泻应及时补充口服补液盐防止脱水",
                "chunk_index": 0,
                "metadata": {"category": "腹泻", "age_range": "1-3岁"},
            },
        ]
        # mock embeddings (1024维随机向量)
        import numpy as np
        np.random.seed(42)
        embeddings = np.random.random((2, 1024)).tolist()

        ids = self.store.insert(chunks, embeddings)
        self.store.flush()
        assert len(ids) == 2

        # 用第一个embedding检索
        results = self.store.search(embeddings[0], top_k=2)
        assert len(results) > 0
        assert "content" in results[0]
        assert "score" in results[0]
        # 第一条应该是最相似的
        assert results[0]["score"] >= results[-1]["score"]

    def test_search_with_filter(self):
        """带元数据过滤的检索"""
        chunks = [
            {"content": "新生儿护理要点", "chunk_index": 0,
             "metadata": {"age_range": "0-1月"}},
            {"content": "幼儿营养搭配指南", "chunk_index": 0,
             "metadata": {"age_range": "1-3岁"}},
        ]
        import numpy as np
        embeddings = np.random.random((2, 1024)).tolist()

        self.store.insert(chunks, embeddings)
        self.store.flush()

        # 只检索0-1月的内容
        results = self.store.search(
            embeddings[0], top_k=5,
            filter_expr='metadata["age_range"] == "0-1月"'
        )
        for r in results:
            assert r["metadata"]["age_range"] == "0-1月"

    def test_count(self):
        """统计向量数量"""
        chunks = [{"content": f"测试内容{i}", "chunk_index": 0, "metadata": {}} for i in range(3)]
        import numpy as np
        embeddings = np.random.random((3, 1024)).tolist()
        self.store.insert(chunks, embeddings)
        self.store.flush()
        assert self.store.count() == 3
