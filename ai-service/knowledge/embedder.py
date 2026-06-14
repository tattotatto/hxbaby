"""Embedding向量化服务 — 基于 bge-large-zh-v1.5"""
from sentence_transformers import SentenceTransformer
from config import get_settings

_embedder_instance = None


class Embedder:
    """文本向量化器，单例模式"""

    def __init__(self, model_name: str, device: str = "cpu"):
        self.model = SentenceTransformer(model_name, device=device)

    def embed(self, text: str) -> list:
        """单文本向量化"""
        embedding = self.model.encode(text, normalize_embeddings=True)
        return embedding.tolist()

    def embed_batch(self, texts: list) -> list:
        """批量向量化"""
        embeddings = self.model.encode(texts, normalize_embeddings=True)
        return embeddings.tolist()


def get_embedder() -> Embedder:
    """获取全局单例Embedder"""
    global _embedder_instance
    if _embedder_instance is None:
        settings = get_settings()
        _embedder_instance = Embedder(
            model_name=settings.embedding_model,
            device=settings.embedding_device,
        )
    return _embedder_instance
