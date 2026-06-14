"""LLM适配器测试"""
import pytest
from unittest.mock import patch, MagicMock


class TestLLMAdapter:
    @pytest.fixture
    def adapter(self):
        from core.llm_adapter import LLMAdapter
        return LLMAdapter(model="qwen-plus")

    def test_generate_returns_content(self, adapter):
        """generate()返回LLM生成的文本"""
        mock_response = MagicMock()
        mock_response.content = "根据您的描述，宝宝可能是消化不良..."

        mock_client = MagicMock()
        mock_client.invoke.return_value = mock_response

        with patch.object(adapter, '_get_client', return_value=mock_client):
            result = adapter.generate([
                {"role": "system", "content": "你是儿科健康顾问"},
                {"role": "user", "content": "宝宝拉奶瓣怎么办？"},
            ])
            assert "消化不良" in result

    def test_stream_generate_yields_tokens(self, adapter):
        """generate_stream()逐token返回"""
        mock_chunks = [
            MagicMock(content="根据"),
            MagicMock(content="您的"),
            MagicMock(content="描述"),
        ]
        mock_client = MagicMock()
        mock_client.stream.return_value = mock_chunks

        with patch.object(adapter, '_get_client', return_value=mock_client):
            tokens = list(adapter.generate_stream([
                {"role": "user", "content": "test"},
            ]))
            assert tokens == ["根据", "您的", "描述"]

    def test_get_available_models(self):
        """get_available_models()返回所有可用模型"""
        from core.llm_adapter import LLMAdapter
        models = LLMAdapter.get_available_models()
        assert "qwen-plus" in models
        assert "deepseek-chat" in models

    def test_model_switching(self):
        """切换模型保持接口一致"""
        from core.llm_adapter import LLMAdapter
        a = LLMAdapter(model="qwen-plus")
        b = LLMAdapter(model="deepseek-chat")
        assert a.model_name == "qwen-plus"
        assert b.model_name == "deepseek-chat"
        # 接口一致，都还没初始化
        assert a._client is None
        assert b._client is None

    def test_lazy_init(self, adapter):
        """懒初始化：构造时不创建client"""
        assert adapter._client is None
        # 调用时才初始化
        mock_client = MagicMock()
        mock_client.invoke.return_value = MagicMock(content="test")
        with patch.object(adapter, '_get_client', return_value=mock_client):
            adapter.generate([{"role": "user", "content": "hi"}])
