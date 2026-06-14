"""症状分析引擎 — LLM驱动的儿童症状评估"""
from typing import List, Dict, Optional
from core.llm_adapter import LLMAdapter


class SymptomAnalyzer:
    """基于知识库和LLM的儿童症状分析"""

    def __init__(self):
        self.llm = LLMAdapter(temperature=0.2)  # 低温度，更确定性

    def analyze(self, user_query: str, knowledge_chunks: List[Dict]) -> Optional[Dict]:
        """
        分析用户描述的症状，结合知识库给出评估。

        Returns:
            {
                "symptoms": ["症状1", ...],
                "possible_conditions": ["可能情况1", ...],
                "risk_level": "low/medium/high/emergency",
                "suggestion_type": "home_care/consult_doctor/emergency",
                "home_care_tips": ["建议1", ...],
                "when_to_see_doctor": "就医指征",
                "analysis": "综合分析"
            }
        """
        if not knowledge_chunks:
            return self._basic_analysis(user_query)

        context = "\n".join([c["content"][:300] for c in knowledge_chunks[:3]])

        prompt = f"""你是一个儿科症状分析专家。请根据知识库内容和用户描述进行症状分析。

知识库参考：
{context}

用户描述：{user_query}

请以JSON格式回复（只返回JSON，不要有其他内容）：
{{
    "symptoms": ["识别到的症状"],
    "possible_conditions": ["可能的健康情况"],
    "risk_level": "low/medium/high/emergency",
    "suggestion_type": "home_care/consult_doctor/emergency",
    "home_care_tips": ["家庭护理建议"],
    "when_to_see_doctor": "需要就医的指征",
    "analysis": "综合分析(50字以内)"
}}"""

        try:
            response = self.llm.generate([{"role": "user", "content": prompt}])
            # 提取JSON（LLM可能在前后加额外文字）
            import json
            json_start = response.find("{")
            json_end = response.rfind("}") + 1
            if json_start >= 0 and json_end > json_start:
                return json.loads(response[json_start:json_end])
        except Exception:
            pass

        return self._basic_analysis(user_query)

    def _basic_analysis(self, query: str) -> Dict:
        """无知识库时的兜底分析"""
        return {
            "symptoms": [],
            "possible_conditions": [],
            "risk_level": "low",
            "suggestion_type": "consult_doctor",
            "home_care_tips": ["建议咨询专业医生获取准确诊断"],
            "when_to_see_doctor": "如症状持续或加重，请及时就医",
            "analysis": "暂无足够信息进行详细分析",
        }
