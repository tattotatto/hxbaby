"""通用文本生成API — 供Go AI Bridge调用"""
from fastapi import APIRouter
from pydantic import BaseModel, Field
from core.llm_adapter import LLMAdapter

router = APIRouter()

# Singleton — default adapter (temperature=0.7 for creative tasks)
_default_adapter = None


def get_adapter(model: str | None = None, max_tokens: int = 500) -> LLMAdapter:
    """获取LLMAdapter实例。使用默认参数时复用单例，自定义参数时新建。"""
    global _default_adapter
    if model is None and max_tokens == 500:
        if _default_adapter is None:
            _default_adapter = LLMAdapter(temperature=0.7, max_tokens=500)
        return _default_adapter
    return LLMAdapter(model=model, temperature=0.7, max_tokens=max_tokens)


class GenerateRequest(BaseModel):
    prompt: str = Field(..., description="提示词/用户输入")
    model: str | None = Field(None, description="模型名称(qwen-plus/deepseek-chat)")
    max_tokens: int = Field(500, ge=1, le=4096, description="最大生成token数")
    parameters: dict | None = Field(None, description="额外参数(预留)")


class GenerateResponse(BaseModel):
    content: str = Field(..., description="生成的文本")
    tokens_used: int = Field(0, description="消耗的token数")
    model: str = Field("", description="使用的模型名称")


@router.post("/generate", response_model=GenerateResponse)
async def generate_text(req: GenerateRequest):
    """通用文本生成（非流式），供Go AI Bridge调用"""
    llm = get_adapter(model=req.model, max_tokens=req.max_tokens)

    # LLMAdapter.generate() 接受 List[Dict] 格式的 messages
    messages = [{"role": "user", "content": req.prompt}]
    content = llm.generate(messages)

    return GenerateResponse(
        content=content,
        tokens_used=0,
        model=llm.model_name,
    )
