"""Embedding向量化服务 — 基于 bge-large-zh-v1.5"""
import os
import logging
from sentence_transformers import SentenceTransformer
from huggingface_hub import snapshot_download
from config import get_settings

logger = logging.getLogger(__name__)
_embedder_instance = None


class Embedder:
    """文本向量化器，单例模式"""

    def __init__(self, model_name: str, device: str = "cpu"):
        # 先从 HuggingFace Hub 下载模型到本地缓存，避免 sentence-transformers 的名称映射问题
        try:
            local_path = snapshot_download(
                repo_id=model_name,
                cache_dir=os.environ.get("SENTENCE_TRANSFORMERS_HOME", None),
                max_workers=2,
                tqdm_class=None,  # disable progress bars in logs
                resume_download=True,
            )
            logger.info(f"Model downloaded to: {local_path}")
            self.model = SentenceTransformer(local_path, device=device)
        except Exception as e:
            logger.warning(f"snapshot_download failed for {model_name}, trying direct load: {e}")
            try:
                self.model = SentenceTransformer(model_name, device=device)
            except Exception as e2:
                raise RuntimeError(f"Failed to load embedding model {model_name}: {e2}")

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
