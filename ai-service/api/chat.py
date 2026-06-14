"""AI对话API — SSE流式核心接口"""
import json
import uuid
import time
from fastapi import APIRouter
from fastapi.responses import StreamingResponse
from models.chat import ChatRequest, ChatResponse, ProductRecommendation, SourceInfo
from core.rag_pipeline import RAGPipeline
from core.llm_adapter import LLMAdapter
from core.symptom_analyzer import SymptomAnalyzer

router = APIRouter()

# 单例
_pipeline = None
_llm = None
_analyzer = None


def get_pipeline() -> RAGPipeline:
    global _pipeline
    if _pipeline is None:
        _pipeline = RAGPipeline()
    return _pipeline


def get_llm() -> LLMAdapter:
    global _llm
    if _llm is None:
        _llm = LLMAdapter(temperature=0.4)
    return _llm


def get_analyzer() -> SymptomAnalyzer:
    global _analyzer
    if _analyzer is None:
        _analyzer = SymptomAnalyzer()
    return _analyzer


SYSTEM_PROMPT = """你是「儿童大健康」AI顾问，基于权威儿科医学知识库为家长提供专业建议。

重要原则：
1. 建议基于知识库中的专业资料，非个人意见
2. 涉及紧急症状（高烧不退、呼吸困难、意识不清等），必须建议立即就医
3. 产品推荐基于专业匹配，非商业推销
4. 回复末尾注明引用的知识来源
5. 始终提醒：AI建议仅供参考，不能替代医生诊断

回复格式建议：
- 先理解问题，表达共情
- 基于知识库给出专业分析
- 给出具体可操作的建议
- 如有必要，推荐匹配产品并说明理由
- 末尾标注来源"""


@router.post("/chat")
async def chat(request: ChatRequest):
    pipeline = get_pipeline()
    llm = get_llm()
    analyzer = get_analyzer()

    # Step 1: 查询重写
    rewritten = pipeline.rewrite_query(request.message)

    # Step 2: 检索（共享知识库 + 年龄过滤）
    knowledge_chunks = pipeline.retrieve(
        rewritten,
        age_months=request.child_age_months,
    )

    # Step 3: 重排序
    ranked = pipeline.rerank(request.message, knowledge_chunks)

    # Step 4: 产品匹配
    products = pipeline.search_products(
        request.message, request.product_chunks or []
    )

    # Step 5: 症状分析
    assessment = analyzer.analyze(request.message, ranked) if ranked else None

    # Step 6: 构建消息
    context = pipeline.build_context(ranked) if ranked else ""

    messages = [{"role": "system", "content": SYSTEM_PROMPT}]
    if request.history:
        messages.extend(request.history[-10:])

    user_msg = request.message
    if context:
        user_msg = f"""知识库参考资料：
{context}

产品推荐参考：
{json.dumps([{"name": p.get("name", ""), "desc": p.get("description", ""), "score": p.get("score", 0)} for p in products[:3]], ensure_ascii=False)}

用户问题：{request.message}

请基于以上参考资料回答，注明来源，末尾给出风险评估。"""

    messages.append({"role": "user", "content": user_msg})

    msg_id = f"msg_{uuid.uuid4().hex[:12]}"

    async def generate():
        # 阶段1: 检索完成
        yield sse_event("thinking", {"stage": "retrieved", "count": len(ranked)})
        time.sleep(0.05)

        # 阶段2: 分析
        yield sse_event("thinking", {"stage": "analyzing"})
        time.sleep(0.05)

        # 阶段3: Token流
        full = ""
        for token in llm.generate_stream(messages):
            full += token
            yield sse_event("token", {"content": token})

        # 阶段4: 完成
        sources = [
            SourceInfo(
                doc_title=c.get("metadata", {}).get("source", ""),
                chunk_index=c.get("chunk_index"),
            )
            for c in ranked[:3]
        ]
        recs = [
            ProductRecommendation(
                product_id=p.get("product_id", ""),
                name=p.get("name", ""),
                reason=p.get("description", p.get("reason", "")),
                score=p.get("score", 0),
            )
            for p in products[:3]
        ]
        done = ChatResponse(
            message_id=msg_id,
            content=full,
            assessment=assessment,
            recommendations=recs,
            sources=sources,
        )
        yield sse_event("done", done.model_dump())

    return StreamingResponse(
        generate(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",
        },
    )


def sse_event(event: str, data) -> str:
    return f"event: {event}\ndata: {json.dumps(data, ensure_ascii=False)}\n\n"
