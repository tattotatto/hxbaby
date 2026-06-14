package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AIBridge struct {
	baseURL    string
	httpClient *http.Client
}

type AIGenerateRequest struct {
	Prompt     string                 `json:"prompt"`
	Model      string                 `json:"model,omitempty"`
	MaxTokens  int                    `json:"max_tokens,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type AIGenerateResponse struct {
	Content string `json:"content"`
	Tokens  int    `json:"tokens_used"`
	Model   string `json:"model"`
}

func NewAIBridge(baseURL string) *AIBridge {
	return &AIBridge{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// generateText calls the Python AI service's generic text generation endpoint
func (b *AIBridge) generateText(prompt string, maxTokens int) (*AIGenerateResponse, error) {
	if maxTokens <= 0 {
		maxTokens = 500
	}
	req := AIGenerateRequest{
		Prompt:    prompt,
		MaxTokens: maxTokens,
	}

	body, _ := json.Marshal(req)
	resp, err := b.httpClient.Post(
		b.baseURL+"/ai/generate",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("AI服务不可用: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI服务返回错误 [%d]: %s", resp.StatusCode, string(bodyBytes))
	}

	var result AIGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}
	return &result, nil
}

// GenerateArticle generates an article based on topic and category
func (b *AIBridge) GenerateArticle(topic, category string) (*AIGenerateResponse, error) {
	prompt := fmt.Sprintf(`你是一个专业的儿童健康科普作者。
请根据以下主题写一篇公众号文章（约800字）：
主题：%s
分类：%s
要求：专业易懂、适合宝妈阅读、结构清晰有小标题。`, topic, category)
	return b.generateText(prompt, 2000)
}

// GenerateSummary generates a summary for an article
func (b *AIBridge) GenerateSummary(article string) (*AIGenerateResponse, error) {
	prompt := fmt.Sprintf("请为以下文章生成一段约100字的摘要：\n\n%s", article)
	return b.generateText(prompt, 200)
}

// GenerateActivityCopy generates marketing copy for an activity
func (b *AIBridge) GenerateActivityCopy(title, description string) (*AIGenerateResponse, error) {
	prompt := fmt.Sprintf(`你是一个活动策划专家。请为以下活动生成营销文案：
活动标题：%s
活动描述：%s
请生成：
1. 朋友圈推广文案（150字以内）
2. 活动详情页文案（300字）
3. 群发通知文案（100字）`, title, description)
	return b.generateText(prompt, 1500)
}

// GenerateSellingPoints extracts core selling points for a product
func (b *AIBridge) GenerateSellingPoints(productName, productDesc string) (*AIGenerateResponse, error) {
	prompt := fmt.Sprintf(`你是一个母婴产品营销专家。请为以下产品提炼核心卖点：
产品名称：%s
产品描述：%s
要求：3-5个卖点，每个卖点20字以内，突出对宝宝/宝妈的价值。`, productName, productDesc)
	return b.generateText(prompt, 500)
}

// GenerateActivityReport generates an activity post-mortem report
func (b *AIBridge) GenerateActivityReport(activityName string, stats map[string]interface{}) (*AIGenerateResponse, error) {
	prompt := fmt.Sprintf(`你是一个数据分析专家。请根据以下活动数据生成复盘报告：
活动名称：%s
活动数据：%v
请分析：参与情况、转化效果、亮点与不足、改进建议。`, activityName, stats)
	return b.generateText(prompt, 1500)
}
