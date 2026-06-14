"""Embedding服务测试"""
import pytest
import numpy as np


class TestEmbedder:
    @pytest.fixture(autouse=True)
    def setup(self):
        """确保embedder加载"""
        from knowledge.embedder import get_embedder
        self.embedder = get_embedder()

    def test_embed_single_text_returns_vector(self):
        """embed()返回正确维度的向量"""
        text = "婴儿消化系统发育指南"
        vector = self.embedder.embed(text)
        assert isinstance(vector, list)
        assert len(vector) > 0
        assert all(isinstance(v, float) for v in vector)

    def test_embed_batch_returns_list_of_vectors(self):
        """embed_batch()返回批量向量"""
        texts = ["婴儿喂养指南", "儿童发育标准", "常见疾病预防"]
        vectors = self.embedder.embed_batch(texts)
        assert len(vectors) == 3
        assert all(len(v) > 0 for v in vectors)
        # 所有向量形状一致
        dims = [len(v) for v in vectors]
        assert len(set(dims)) == 1

    def test_embedding_dimension_matches_config(self):
        """向量维度与配置一致"""
        from config import get_settings
        vector = self.embedder.embed("test")
        assert len(vector) == get_settings().embedding_dim

    def test_similar_texts_have_higher_similarity(self):
        """语义相似文本的余弦相似度高于不相关文本"""
        v1 = self.embedder.embed("宝宝消化不好怎么办")
        v2 = self.embedder.embed("婴儿消化不良的处理方法")
        v3 = self.embedder.embed("今天天气很好适合出去玩")

        sim_12 = _cosine_sim(v1, v2)
        sim_13 = _cosine_sim(v1, v3)

        assert sim_12 > sim_13, (
            f"相似文本相似度({sim_12:.4f})应高于不相关文本({sim_13:.4f})"
        )

    def test_singleton_returns_same_instance(self):
        """get_embedder()返回同一实例"""
        from knowledge.embedder import get_embedder
        e1 = get_embedder()
        e2 = get_embedder()
        assert e1 is e2


def _cosine_sim(a, b):
    """计算余弦相似度"""
    a, b = np.array(a), np.array(b)
    return float(np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b)))
