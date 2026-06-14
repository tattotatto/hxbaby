package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

type ChatPayload struct {
	Message        string                   `json:"message"`
	ConversationID string                   `json:"conversation_id,omitempty"`
	ChildID        string                   `json:"child_id,omitempty"`
	ChildAgeMonths int                      `json:"child_age_months,omitempty"`
	TenantID       string                   `json:"tenant_id"`
	ProductChunks  []map[string]interface{} `json:"product_chunks,omitempty"`
	History        []map[string]interface{} `json:"history,omitempty"`
}

func NewAIClient(baseURL string) *AIClient {
	return &AIClient{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{},
	}
}

func (c *AIClient) ChatStream(payload ChatPayload, writer io.Writer) error {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", c.BaseURL+"/ai/chat", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("AI service request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("AI service unavailable: %w", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		writer.Write([]byte(line + "\n"))
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}
	}
	return scanner.Err()
}
