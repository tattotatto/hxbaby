"""文档加载器 - 支持 PDF/DOCX/TXT/Markdown/CSV/Excel"""
import os
from pathlib import Path
from typing import List, Dict

SUPPORTED_FORMATS = {".pdf", ".docx", ".txt", ".md", ".csv", ".xlsx"}


class DocumentLoadError(Exception):
    """文档加载异常"""
    pass


def load_document(file_path: str) -> List[Dict]:
    """
    加载文档，统一返回 [{content, metadata}] 格式。

    Args:
        file_path: 文档文件路径

    Returns:
        [{"content": "文档内容", "metadata": {"source": "...", "format": "..."}}]

    Raises:
        DocumentLoadError: 文件不存在或格式不支持
    """
    if not os.path.exists(file_path):
        raise DocumentLoadError(f"File not found: {file_path}")

    ext = Path(file_path).suffix.lower()
    if ext not in SUPPORTED_FORMATS:
        raise DocumentLoadError(
            f"Unsupported format: {ext}. Supported: {SUPPORTED_FORMATS}"
        )

    if ext in (".txt", ".md"):
        return _load_text(file_path)
    elif ext == ".pdf":
        return _load_pdf(file_path)
    elif ext == ".docx":
        return _load_docx(file_path)
    elif ext in (".csv", ".xlsx"):
        return _load_table(file_path)

    return []


def _load_text(file_path: str) -> List[Dict]:
    """加载纯文本/Markdown文件"""
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()
    return [{
        "content": content,
        "metadata": {"source": file_path, "format": "text"},
    }]


def _load_pdf(file_path: str) -> List[Dict]:
    """加载PDF文件 (使用unstructured库)"""
    try:
        from unstructured.partition.pdf import partition_pdf
        elements = partition_pdf(file_path)
        content = "\n".join([str(el) for el in elements])
        return [{
            "content": content,
            "metadata": {"source": file_path, "format": "pdf"},
        }]
    except ImportError:
        raise DocumentLoadError(
            "PDF support requires 'unstructured' library. Run: pip install unstructured"
        )
    except Exception as e:
        raise DocumentLoadError(f"Failed to load PDF: {e}")


def _load_docx(file_path: str) -> List[Dict]:
    """加载Word文档"""
    try:
        from unstructured.partition.docx import partition_docx
        elements = partition_docx(file_path)
        content = "\n".join([str(el) for el in elements])
        return [{
            "content": content,
            "metadata": {"source": file_path, "format": "docx"},
        }]
    except ImportError:
        raise DocumentLoadError(
            "DOCX support requires 'unstructured' library. Run: pip install unstructured"
        )
    except Exception as e:
        raise DocumentLoadError(f"Failed to load DOCX: {e}")


def _load_table(file_path: str) -> List[Dict]:
    """加载CSV/Excel表格，转为Markdown格式"""
    ext = Path(file_path).suffix.lower()
    try:
        import pandas as pd
        if ext == ".csv":
            df = pd.read_csv(file_path, encoding="utf-8")
        else:
            df = pd.read_excel(file_path)
        content = df.to_markdown(index=False)
        return [{
            "content": content,
            "metadata": {
                "source": file_path,
                "format": "table",
                "rows": len(df),
                "columns": list(df.columns),
            },
        }]
    except ImportError:
        raise DocumentLoadError(
            "Table support requires 'pandas' and 'tabulate'. Run: pip install pandas tabulate openpyxl"
        )
    except Exception as e:
        raise DocumentLoadError(f"Failed to load table: {e}")
