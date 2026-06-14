"""儿童大健康AI服务 — FastAPI入口"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
from config import get_settings

settings = get_settings()


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup: 初始化 Embedding 模型和 Milvus 连接
    import logging
    logger = logging.getLogger(__name__)
    logger.info("AI Service starting up...")
    try:
        from knowledge.embedder import get_embedder
        get_embedder()
        logger.info(f"Embedding model loaded: {settings.embedding_model}")
    except Exception as e:
        logger.warning(f"Embedding model not initialized (will lazy-load): {e}")

    try:
        from knowledge.vector_store import get_vector_store
        get_vector_store()
        logger.info(f"Milvus connected: {settings.milvus_host}:{settings.milvus_port}")
    except Exception as e:
        logger.warning(f"Milvus not connected (will lazy-load): {e}")

    yield
    logger.info("AI Service shutting down...")


app = FastAPI(
    title="儿童大健康AI服务",
    description="基于RAG的儿童健康智能问答系统",
    version="0.1.0",
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 注册路由
from api.health import router as health_router
app.include_router(health_router, prefix="/ai", tags=["health"])

from api.chat import router as chat_router
app.include_router(chat_router, prefix="/ai", tags=["chat"])

from api.generate import router as generate_router
app.include_router(generate_router, prefix="/ai", tags=["generate"])


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "main:app",
        host=settings.ai_service_host,
        port=settings.ai_service_port,
        reload=settings.debug,
    )
