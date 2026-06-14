"""RAG检索管道测试"""
import pytest
from unittest.mock import patch, MagicMock


class TestRAGPipeline:
    @pytest.fixture
    def pipeline(self):
        """创建pipeline，mock底层依赖"""
        with patch('core.rag_pipeline.get_vector_store') as mock_vs, \
             patch('core.rag_pipeline.get_embedder') as mock_emb:
            mock_store = MagicMock()
            mock_emb_instance = MagicMock()
            mock_emb_instance.embed.return_value = [0.1] * 1024
            mock_vs.return_value = mock_store
            mock_emb.return_value = mock_emb_instance

            from core.rag_pipeline import RAGPipeline
            p = RAGPipeline()
            p.vector_store = mock_store
            p.embedder = mock_emb_instance
            yield p

    def test_retrieve_calls_vector_search(self, pipeline):
        mock_hits = [
            {"content": "婴儿发热处理指南", "score": 0.92, "metadata": {"category": "发热"}},
            {"content": "儿童体温测量方法", "score": 0.85, "metadata": {"category": "护理"}},
        ]
        pipeline.vector_store.search.return_value = mock_hits

        results = pipeline.retrieve("宝宝发烧", top_k=5)
        assert len(results) == 2
        assert results[0]["content"] == "婴儿发热处理指南"

    def test_retrieve_with_age_filter(self, pipeline):
        pipeline.vector_store.search.return_value = []
        pipeline.retrieve("喂养", age_months=6)
        call_args = pipeline.vector_store.search.call_args
        assert call_args[1]["filter_expr"] is not None

    def test_rewrite_query_fallback(self, pipeline):
        assert pipeline.rewrite_query("宝宝拉肚子") is not None

    def test_build_context_formats(self, pipeline):
        chunks = [
            {"content": "知识点1", "metadata": {"source": "doc1.pdf"}},
            {"content": "知识点2", "metadata": {"source": "doc2.pdf"}},
        ]
        ctx = pipeline.build_context(chunks)
        assert "知识点1" in ctx
        assert "doc1.pdf" in ctx

    def test_rerank_returns_top_k(self, pipeline):
        chunks = [{"content": f"内容{i}", "score": 1.0 - i * 0.1, "metadata": {}} for i in range(10)]
        ranked = pipeline.rerank("测试", chunks, top_k=3)
        assert len(ranked) <= 3

    def test_search_products_empty(self, pipeline):
        assert pipeline.search_products("test", []) == []
