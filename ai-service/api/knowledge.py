"""知识库管理 API — 文档上传/列表/删除/统计"""
import os
import tempfile
import uuid
from datetime import datetime
from pathlib import Path

from fastapi import APIRouter, File, UploadFile, HTTPException
from pydantic import BaseModel

router = APIRouter(tags=["knowledge"])

MAX_FILE_SIZE = 50 * 1024 * 1024  # 50MB


class UploadResponse(BaseModel):
    source: str
    chunks: int
    format: str
    status: str


class DocumentItem(BaseModel):
    source: str
    format: str
    chunk_count: int
    indexed_at: str


class DocumentListResponse(BaseModel):
    documents: list
    total: int


class DeleteResponse(BaseModel):
    source: str
    deleted_chunks: int


class StatsResponse(BaseModel):
    total_chunks: int
    total_documents: int
    collection_name: str


@router.post("/upload", response_model=UploadResponse)
async def upload_document(file: UploadFile = File(...)):
    """上传文档并入库 (PDF/DOCX/TXT/MD/CSV/XLSX)"""
    from knowledge.loader import load_document, DocumentLoadError, SUPPORTED_FORMATS
    from knowledge.chunker import chunk_document
    from knowledge.embedder import get_embedder
    from knowledge.vector_store import get_vector_store
    from config import get_settings
    import logging

    logger = logging.getLogger(__name__)
    settings = get_settings()

    # 1. 校验格式
    ext = Path(file.filename or "unknown").suffix.lower()
    if ext not in SUPPORTED_FORMATS:
        raise HTTPException(
            status_code=400,
            detail=f"不支持的格式: {ext}。支持: {', '.join(sorted(SUPPORTED_FORMATS))}"
        )

    # 2. 保存临时文件
    tmp_path = None
    try:
        suffix = ext if ext else ""
        with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
            content_bytes = await file.read()
            if len(content_bytes) > MAX_FILE_SIZE:
                raise HTTPException(status_code=400, detail=f"文件超过 {MAX_FILE_SIZE // 1024 // 1024}MB 限制")
            tmp.write(content_bytes)
            tmp_path = tmp.name

        # 3. 文档加载
        docs = load_document(tmp_path)
        if not docs:
            raise HTTPException(status_code=400, detail="文档内容为空")

        # 4. 分块
        source_name = file.filename or f"upload_{uuid.uuid4().hex[:8]}"
        all_chunks = []
        for doc in docs:
            # 将临时路径替换为原始文件名
            doc["metadata"]["source"] = source_name
            chunks = chunk_document(
                doc["content"],
                chunk_size=settings.chunk_size,
                overlap=settings.chunk_overlap,
                metadata=doc["metadata"],
            )
            all_chunks.extend(chunks)

        if not all_chunks:
            raise HTTPException(status_code=400, detail="文档分块后无有效内容")

        # 5. 向量化
        embedder = get_embedder()
        texts = [c["content"] for c in all_chunks]
        embeddings = embedder.embed_batch(texts)

        # 6. 写入 Milvus（添加时间戳）
        for c in all_chunks:
            c["metadata"]["indexed_at"] = datetime.utcnow().isoformat()

        vector_store = get_vector_store()
        pk_ids = vector_store.insert(all_chunks, embeddings)
        vector_store.flush()

        logger.info(f"Document indexed: {source_name}, chunks={len(all_chunks)}, ids={pk_ids}")

        return UploadResponse(
            source=source_name,
            chunks=len(all_chunks),
            format=ext.lstrip("."),
            status="indexed",
        )

    except HTTPException:
        raise
    except DocumentLoadError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.error(f"Failed to index document: {e}")
        raise HTTPException(status_code=500, detail=f"文档处理失败: {e}")
    finally:
        if tmp_path and os.path.exists(tmp_path):
            os.unlink(tmp_path)


@router.get("/documents", response_model=DocumentListResponse)
async def list_documents():
    """列出所有已入库文档（按 source 聚合）"""
    from knowledge.vector_store import get_vector_store
    from pymilvus import Collection

    try:
        vector_store = get_vector_store()
        vector_store._connect()
        collection = vector_store.collection

        # 查询所有实体，按 metadata.source 去重
        collection.load()
        # 查询全部 — 获取 content, metadata 字段（分批，避免超 Milvus 限制）
        count = collection.num_entities
        if count == 0:
            return DocumentListResponse(documents=[], total=0)
        all_results = []
        batch_size = 8000
        for offset in range(0, count, batch_size):
            batch = collection.query(
                expr="id >= 0",
                output_fields=["content", "chunk_index", "metadata"],
                limit=batch_size,
                offset=offset,
            )
            all_results.extend(batch)

        # 按 source 聚合
        doc_map = {}
        for r in results:
            meta = r.get("metadata", {})
            source = meta.get("source", "unknown")
            fmt = meta.get("format", "unknown")
            indexed_at = meta.get("indexed_at", "")
            if source not in doc_map:
                doc_map[source] = {
                    "source": source,
                    "format": fmt,
                    "chunk_count": 0,
                    "indexed_at": indexed_at,
                }
            doc_map[source]["chunk_count"] += 1
            # 保留最新的 indexed_at
            if indexed_at and indexed_at > doc_map[source]["indexed_at"]:
                doc_map[source]["indexed_at"] = indexed_at

        documents = sorted(doc_map.values(), key=lambda d: d["indexed_at"], reverse=True)
        return DocumentListResponse(documents=documents, total=len(documents))

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"查询失败: {e}")


@router.delete("/documents/{source:path}", response_model=DeleteResponse)
async def delete_document(source: str):
    """按 source 删除文档的所有 chunks"""
    from knowledge.vector_store import get_vector_store

    try:
        vector_store = get_vector_store()
        vector_store._connect()
        collection = vector_store.collection

        # 先查询该 source 有多少 chunks（分批查询）
        collection.load()
        count = collection.num_entities
        all_ids = []
        batch_size = 8000
        source_expr = f"metadata['source'] == '{source}'"
        for offset in range(0, count, batch_size):
            batch = collection.query(
                expr=source_expr,
                output_fields=["id"],
                limit=batch_size,
                offset=offset,
            )
            all_ids.extend([r["id"] for r in batch])
        chunk_count = len(all_ids)

        if chunk_count == 0:
            raise HTTPException(status_code=404, detail=f"文档不存在: {source}")

        # 逐条删除
        ids_to_delete = [str(id_) for id_ in all_ids]
        # 分批删除（Milvus 限制表达式长度）
        for i in range(0, len(ids_to_delete), 500):
            batch_ids = ids_to_delete[i:i+500]
            collection.delete(f"id in [{','.join(batch_ids)}]")
        collection.flush()

        return DeleteResponse(source=source, deleted_chunks=chunk_count)

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"删除失败: {e}")


@router.get("/stats", response_model=StatsResponse)
async def get_stats():
    """获取知识库统计信息"""
    from knowledge.vector_store import get_vector_store

    try:
        vector_store = get_vector_store()
        total_chunks = vector_store.count()

        # 统计去重文档数
        doc_count = 0
        if total_chunks > 0:
            vector_store._connect()
            collection = vector_store.collection
            collection.load()
            all_results = []
            batch_size = 8000
            for offset in range(0, total_chunks, batch_size):
                batch = collection.query(
                    expr="id >= 0",
                    output_fields=["metadata"],
                    limit=batch_size,
                    offset=offset,
                )
                all_results.extend(batch)
            sources = set()
            for r in all_results:
                meta = r.get("metadata", {})
                src = meta.get("source", "")
                if src:
                    sources.add(src)
            doc_count = len(sources)

        return StatsResponse(
            total_chunks=total_chunks,
            total_documents=doc_count,
            collection_name=vector_store.collection_name,
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"查询失败: {e}")
