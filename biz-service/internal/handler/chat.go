package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/service"
)

type ChatHandler struct {
	aiClient *service.AIClient
}

func NewChatHandler(aiClient *service.AIClient) *ChatHandler {
	return &ChatHandler{aiClient: aiClient}
}

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
	ChildID string `json:"child_id"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	conversationID := c.Param("id")

	payload := service.ChatPayload{
		Message:        req.Message,
		ConversationID: conversationID,
		ChildID:        req.ChildID,
		TenantID:       tenantID,
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	if err := h.aiClient.ChatStream(payload, c.Writer); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
}
