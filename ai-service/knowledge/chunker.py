"""文档分块器 - 递归语义分块，支持重叠"""
from typing import List, Dict, Optional


def chunk_document(
    text: str,
    chunk_size: int = 800,
    overlap: int = 120,
    metadata: Optional[Dict] = None,
    separators: List[str] = None,
) -> List[Dict]:
    """
    将文档按语义分块，优先按段落→句子→字符分割。

    Args:
        text: 原始文本
        chunk_size: 每块最大字符数
        overlap: 相邻块重叠字符数
        metadata: 附加元数据
        separators: 分隔符优先级列表

    Returns:
        [{"content": "...", "chunk_index": 0, "metadata": {...}}, ...]
    """
    if separators is None:
        separators = ["\n\n", "\n", "。", ".", "；", ";", "，", ",", " "]

    chunks_text = _split_text(text, separators, chunk_size, overlap)

    result = []
    for i, chunk_text in enumerate(chunks_text):
        chunk_meta = (metadata or {}).copy()
        chunk_meta["chunk_size"] = len(chunk_text)
        result.append({
            "content": chunk_text,
            "chunk_index": i,
            "metadata": chunk_meta,
        })

    return result


def _split_text(text: str, separators: List[str], chunk_size: int, overlap: int) -> List[str]:
    """将长文本按分隔符优先级切分"""
    if len(text) <= chunk_size:
        return [text.strip()] if text.strip() else []

    # 尝试用第一个分隔符切分
    sep = separators[0] if separators else None

    if sep and sep in text:
        parts = text.split(sep)
        chunks = []
        current = ""
        for part in parts:
            candidate = (current + sep + part) if current else part
            if len(candidate) <= chunk_size:
                current = candidate
            else:
                if current.strip():
                    chunks.append(current.strip())
                # 单个part仍太长，递归用下一个分隔符处理
                if len(part) > chunk_size:
                    sub_chunks = _split_text(
                        part, separators[1:], chunk_size, overlap
                    ) if len(separators) > 1 else _force_split(part, chunk_size, overlap)
                    chunks.extend(sub_chunks)
                    current = ""
                else:
                    current = (current[-overlap:] + part) if current and overlap > 0 else part
        if current.strip():
            chunks.append(current.strip())
        return chunks

    # 当前分隔符不匹配，用下一个
    if len(separators) > 1:
        return _split_text(text, separators[1:], chunk_size, overlap)

    # 所有分隔符都不匹配，强制按长度切
    return _force_split(text, chunk_size, overlap)


def _force_split(text: str, chunk_size: int, overlap: int) -> List[str]:
    """强制按字符数切分（无分隔符可用时的兜底方案）"""
    chunks = []
    start = 0
    while start < len(text):
        end = min(start + chunk_size, len(text))
        chunks.append(text[start:end].strip())
        if end >= len(text):
            break
        start = end - overlap
    return [c for c in chunks if c]
