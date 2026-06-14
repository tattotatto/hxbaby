"""对话相关数据模型"""
from pydantic import BaseModel, Field
from typing import Optional, List, Dict


class ChatRequest(BaseModel):
    """AI对话请求"""
    message: str = Field(..., description="用户消息")
    conversation_id: Optional[str] = Field(None, description="对话ID")
    child_id: Optional[str] = Field(None, description="儿童ID")
    child_age_months: Optional[int] = Field(None, description="儿童月龄")
    tenant_id: str = Field(..., description="租户ID")
    product_chunks: Optional[List[Dict]] = Field(
        default_factory=list, description="租户产品库(含embedding)"
    )
    history: Optional[List[Dict]] = Field(
        default_factory=list, description="历史消息(最多10条)"
    )


class SourceInfo(BaseModel):
    """知识来源"""
    doc_title: Optional[str] = None
    chunk_index: Optional[int] = None


class ProductRecommendation(BaseModel):
    """产品推荐"""
    product_id: str
    name: str
    reason: str
    score: float


class ChatResponse(BaseModel):
    """AI对话完整响应"""
    message_id: str
    content: str
    assessment: Optional[Dict] = None
    recommendations: List[ProductRecommendation] = Field(default_factory=list)
    sources: List[SourceInfo] = Field(default_factory=list)
