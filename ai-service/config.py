"""应用配置管理"""
from pydantic_settings import BaseSettings
from functools import lru_cache
import os


class Settings(BaseSettings):
    # 服务
    ai_service_host: str = "0.0.0.0"
    ai_service_port: int = 8001
    debug: bool = False

    # LLM — 主模型
    llm_model: str = "deepseek-v4-pro"
    qwen_api_key: str = ""
    qwen_base_url: str = "https://dashscope.aliyuncs.com/compatible-mode/v1"
    deepseek_api_key: str = ""
    deepseek_base_url: str = "https://api.deepseek.com/v1"

    # Embedding
    embedding_model: str = "BAAI/bge-large-zh-v1.5"
    embedding_device: str = "cpu"
    embedding_dim: int = 1024

    # Reranker
    reranker_model: str = "BAAI/bge-reranker-v2-m3"

    # Milvus
    milvus_host: str = "localhost"
    milvus_port: int = 19530
    milvus_collection_name: str = "knowledge_chunks"

    # 分块
    chunk_size: int = 800
    chunk_overlap: int = 120

    # 检索
    retrieval_top_k: int = 10
    rerank_top_k: int = 5

    # Redis
    redis_url: str = "redis://localhost:6379/0"

    class Config:
        env_file = ".env"
        extra = "allow"


@lru_cache()
def get_settings() -> Settings:
    return Settings()
