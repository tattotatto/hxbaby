#!/usr/bin/env python
"""知识库初始化脚本 — 批量导入文档到向量数据库

用法:
    python scripts/init_knowledge.py --dir ./knowledge_docs/
    python scripts/init_knowledge.py --dir ./knowledge_docs/ --category 营养学
    python scripts/init_knowledge.py --dir ./knowledge_docs/ --dry-run
"""
import argparse
import os
import sys
import time
from pathlib import Path

# 添加 ai-service 到路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "ai-service"))

from knowledge.loader import load_document, DocumentLoadError
from knowledge.chunker import chunk_document
from knowledge.embedder import get_embedder
from knowledge.vector_store import get_vector_store
from config import get_settings


def main():
    parser = argparse.ArgumentParser(description="儿童健康知识库初始化")
    parser.add_argument("--dir", required=True, help="文档目录路径")
    parser.add_argument("--category", default=None, help="文档分类标签")
    parser.add_argument("--age-range", default=None, help="适用年龄段")
    parser.add_argument("--dry-run", action="store_true", help="仅展示处理流程，不实际写入")
    args = parser.parse_args()

    if not os.path.isdir(args.dir):
        print(f"错误: 目录不存在: {args.dir}")
        sys.exit(1)

    settings = get_settings()
    embedder = get_embedder()
    vector_store = get_vector_store()

    # 收集所有支持的文件
    supported_exts = {".pdf", ".docx", ".txt", ".md", ".csv", ".xlsx"}
    files = []
    for root, _, filenames in os.walk(args.dir):
        for fname in filenames:
            if Path(fname).suffix.lower() in supported_exts:
                files.append(os.path.join(root, fname))

    print(f"发现 {len(files)} 个文档待处理\n")

    total_chunks = 0
    total_errors = 0
    start_time = time.time()

    for i, file_path in enumerate(files, 1):
        fname = os.path.basename(file_path)
        print(f"[{i}/{len(files)}] {fname} ... ", end="", flush=True)

        try:
            # Step 1: 加载
            docs = load_document(file_path)

            # Step 2: 分块
            metadata = {"source": fname}
            if args.category:
                metadata["category"] = args.category
            if args.age_range:
                metadata["age_range"] = args.age_range

            all_chunks = []
            for doc in docs:
                doc_meta = {**doc["metadata"], **metadata}
                chunks = chunk_document(
                    doc["content"],
                    chunk_size=settings.chunk_size,
                    overlap=settings.chunk_overlap,
                    metadata=doc_meta,
                )
                all_chunks.extend(chunks)

            if args.dry_run:
                print(f"OK ({len(all_chunks)} chunks) [DRY RUN]")
                total_chunks += len(all_chunks)
                continue

            # Step 3: 向量化
            texts = [c["content"] for c in all_chunks]
            embeddings = embedder.embed_batch(texts)

            # Step 4: 写入Milvus
            vector_store.insert(all_chunks, embeddings)
            total_chunks += len(all_chunks)

            print(f"OK ({len(all_chunks)} chunks)")

        except DocumentLoadError as e:
            print(f"SKIP: {e}")
            total_errors += 1
        except Exception as e:
            print(f"ERROR: {e}")
            total_errors += 1

    if not args.dry_run:
        vector_store.flush()

    elapsed = time.time() - start_time
    print(f"\n{'='*50}")
    print(f"完成! 总计: {total_chunks} chunks, {total_errors} errors")
    print(f"耗时: {elapsed:.1f}s")
    if total_chunks > 0:
        print(f"平均: {elapsed/total_chunks:.2f}s/chunk")


if __name__ == "__main__":
    main()
