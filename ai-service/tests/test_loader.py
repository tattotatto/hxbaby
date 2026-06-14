"""文档加载器测试"""
import pytest
from knowledge.loader import load_document, DocumentLoadError


class TestDocumentLoader:
    def test_load_txt_returns_content(self, tmp_path):
        """加载TXT文件应返回正确内容"""
        txt_path = tmp_path / "test.txt"
        txt_path.write_text("婴儿消化系统发育指南\n第一章 新生儿的消化特点", encoding="utf-8")
        docs = load_document(str(txt_path))
        assert len(docs) == 1
        assert "婴儿消化系统发育指南" in docs[0]["content"]
        assert docs[0]["metadata"]["format"] == "text"

    def test_load_markdown_returns_content(self, tmp_path):
        """加载Markdown文件"""
        md_path = tmp_path / "guide.md"
        md_path.write_text("# 婴儿喂养指南\n\n## 母乳喂养\n母乳是婴儿最好的食物。", encoding="utf-8")
        docs = load_document(str(md_path))
        assert len(docs) == 1
        assert "母乳喂养" in docs[0]["content"]

    def test_load_unsupported_format_raises_error(self, tmp_path):
        """不支持的格式应抛出DocumentLoadError"""
        bad_path = tmp_path / "test.xyz"
        bad_path.write_text("content")
        with pytest.raises(DocumentLoadError, match="Unsupported format"):
            load_document(str(bad_path))

    def test_load_nonexistent_file_raises_error(self):
        """文件不存在应报错"""
        with pytest.raises(DocumentLoadError, match="File not found"):
            load_document("/nonexistent/file.pdf")

    def test_load_empty_file_returns_empty_content(self, tmp_path):
        """空文件返回空内容"""
        txt_path = tmp_path / "empty.txt"
        txt_path.write_text("", encoding="utf-8")
        docs = load_document(str(txt_path))
        assert len(docs) == 1
        assert docs[0]["content"] == ""

    def test_metadata_includes_source_path(self, tmp_path):
        """metadata应包含源文件路径"""
        txt_path = tmp_path / "health.txt"
        txt_path.write_text("测试内容", encoding="utf-8")
        docs = load_document(str(txt_path))
        assert docs[0]["metadata"]["source"] == str(txt_path)
