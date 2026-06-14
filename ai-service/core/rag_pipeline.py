"""RAG检索增强生成管道"""
from typing import List, Dict, Optional
import numpy as np
from knowledge.vector_store import get_vector_store
from knowledge.embedder import get_embedder
from core.llm_adapter import LLMAdapter
from config import get_settings


class RAGPipeline:
    """RAG管道：检索→重排→上下文构建→产品匹配"""

    def __init__(self):
        self.vector_store = get_vector_store()
        self.embedder = get_embedder()
        self.llm = LLMAdapter()
        self.settings = get_settings()

    def rewrite_query(self, query: str) -> str:
        """查询重写：将口语化问题改写为精准的医学检索查询"""
        prompt = f"""你是一个儿童健康领域的查询优化专家。
请将用户的原始问题改写为更适合医学知识检索的查询（保持原意，添加医学术语）。

原始问题：{query}

改写后的检索查询："""
        try:
            result = self.llm.generate([{"role": "user", "content": prompt}])
            return result.strip() or query
        except Exception:
            return query  # LLM不可用时返回原查询

    def retrieve(
        self,
        query: str,
        top_k: int = None,
        age_months: int = None,
        category: str = None,
    ) -> List[Dict]:
        """检索相关知识块，支持年龄段和分类过滤"""
        if top_k is None:
            top_k = self.settings.retrieval_top_k

        # 构建过滤表达式
        filters = []
        if age_months is not None:
            filters.append(_build_age_filter(age_months))
        if category:
            filters.append(f'metadata["category"] == "{category}"')

        filter_expr = " && ".join(filters) if filters else None

        query_vec = self.embedder.embed(query)
        return self.vector_store.search(query_vec, top_k=top_k, filter_expr=filter_expr)

    def rerank(self, query: str, chunks: List[Dict], top_k: int = None) -> List[Dict]:
        """LLM重排序：从候选chunks中选出最相关的"""
        if top_k is None:
            top_k = self.settings.rerank_top_k
        if len(chunks) <= top_k:
            return chunks

        # 用LLM选择最相关的N条
        chunk_texts = "\n---\n".join([
            f"[{i}] {c['content'][:200]}" for i, c in enumerate(chunks)
        ])
        prompt = f"""从以下检索结果中选出与问题最相关的{top_k}条，只返回编号列表(如: 2,0,5)。

问题：{query}

检索结果：
{chunk_texts}

最相关的{top_k}条编号（逗号分隔）："""

        try:
            response = self.llm.generate([{"role": "user", "content": prompt}])
            indices = [
                int(x.strip()) for x in response.split(",")[:top_k]
                if x.strip().isdigit()
            ]
            return [chunks[i] for i in indices if 0 <= i < len(chunks)]
        except Exception:
            return chunks[:top_k]

    def build_context(self, chunks: List[Dict]) -> str:
        """将检索到的chunks组合为LLM上下文"""
        parts = []
        for i, chunk in enumerate(chunks):
            src = chunk.get("metadata", {}).get("source", "未知来源")
            parts.append(f"【参考资料{i+1}】(来源: {src})\n{chunk['content']}")
        return "\n\n---\n\n".join(parts)

    def search_products(
        self, query: str, product_chunks: List[Dict], top_k: int = 3
    ) -> List[Dict]:
        """在租户产品库中匹配相关产品"""
        if not product_chunks:
            return []
        query_vec = self.embedder.embed(query)
        scored = []
        for pc in product_chunks:
            emb = pc.get("embedding")
            if emb:
                sim = float(np.dot(query_vec, emb) / (np.linalg.norm(query_vec) * np.linalg.norm(emb)))
                scored.append({**pc, "score": round(sim, 4)})
        scored.sort(key=lambda x: x.get("score", 0), reverse=True)
        return scored[:top_k]


def _build_age_filter(age_months: int) -> str:
    """根据月龄构建Milvus年龄过滤表达式"""
    if age_months <= 1:
        return 'metadata["age_range"] in ["0-1月", "新生儿"]'
    elif age_months <= 12:
        return 'metadata["age_range"] in ["0-1月", "1-12月", "婴儿期", "0-1岁"]'
    elif age_months <= 36:
        return 'metadata["age_range"] in ["1-3岁", "幼儿期"]'
    else:
        return 'metadata["age_range"] in ["3-6岁", "学龄前", "儿童期"]'
