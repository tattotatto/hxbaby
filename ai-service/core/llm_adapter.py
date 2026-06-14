"""LLM适配器 — 统一封装通义千问/DeepSeek，支持同步和流式调用"""
from typing import List, Dict, Iterator, Optional
from config import get_settings


class LLMAdapter:
    """LLM适配器，封装多模型，提供统一接口（懒初始化）"""

    MODEL_CONFIGS = {
        # 通义千问系列
        "qwen-plus": {
            "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
            "provider": "qwen",
        },
        "qwen-max": {
            "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
            "provider": "qwen",
        },
        # DeepSeek 系列
        "deepseek-chat": {
            "base_url": "https://api.deepseek.com/v1",
            "provider": "deepseek",
        },
        "deepseek-v4-pro": {
            "base_url": "https://api.deepseek.com/v1",
            "provider": "deepseek",
        },
    }

    def __init__(self, model: str = None, temperature: float = 0.3, max_tokens: int = 2048):
        settings = get_settings()
        self.model_name = model or settings.llm_model
        self.temperature = temperature
        self.max_tokens = max_tokens
        self._client = None  # 懒初始化

    def _get_client(self):
        """懒初始化 ChatOpenAI 客户端"""
        if self._client is not None:
            return self._client

        from langchain_openai import ChatOpenAI
        import os

        settings = get_settings()
        config = self.MODEL_CONFIGS.get(self.model_name, self.MODEL_CONFIGS["deepseek-v4-pro"])
        provider = config.get("provider", "deepseek")

        # 根据 provider 选择对应的 API Key
        if provider == "deepseek":
            api_key = settings.deepseek_api_key or os.environ.get("DEEPSEEK_API_KEY", "")
        elif provider == "qwen":
            api_key = settings.qwen_api_key or os.environ.get("QWEN_API_KEY", "")
        else:
            api_key = ""

        self._client = ChatOpenAI(
            model=self.model_name,
            api_key=api_key or "dummy-key",  # 测试用
            base_url=config["base_url"],
            temperature=self.temperature,
            max_tokens=self.max_tokens,
        )
        return self._client

    def generate(self, messages: List[Dict]) -> str:
        """同步生成回复"""
        from langchain_core.messages import HumanMessage, SystemMessage, AIMessage

        lc_messages = []
        for m in messages:
            role = m["role"]
            content = m["content"]
            if role == "system":
                lc_messages.append(SystemMessage(content=content))
            elif role == "user":
                lc_messages.append(HumanMessage(content=content))
            elif role == "assistant":
                lc_messages.append(AIMessage(content=content))

        response = self._get_client().invoke(lc_messages)
        return response.content

    def generate_stream(self, messages: List[Dict]) -> Iterator[str]:
        """流式生成回复，逐token yield"""
        from langchain_core.messages import HumanMessage, SystemMessage, AIMessage

        lc_messages = []
        for m in messages:
            role = m["role"]
            content = m["content"]
            if role == "system":
                lc_messages.append(SystemMessage(content=content))
            elif role == "user":
                lc_messages.append(HumanMessage(content=content))
            elif role == "assistant":
                lc_messages.append(AIMessage(content=content))

        for chunk in self._get_client().stream(lc_messages):
            if chunk.content:
                yield chunk.content

    @classmethod
    def get_available_models(cls) -> List[str]:
        """返回可用模型列表"""
        return list(cls.MODEL_CONFIGS.keys())
