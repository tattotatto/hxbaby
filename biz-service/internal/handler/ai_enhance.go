package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/service"
	"github.com/hxbaby/biz-service/pkg/response"
)

type AIEnhanceHandler struct {
	bridge *service.AIBridge
}

func NewAIEnhanceHandler(bridge *service.AIBridge) *AIEnhanceHandler {
	return &AIEnhanceHandler{bridge: bridge}
}

func (h *AIEnhanceHandler) GenerateArticle(c *gin.Context) {
	var req struct {
		Topic    string `json:"topic" binding:"required"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.bridge.GenerateArticle(req.Topic, req.Category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateSummary(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.bridge.GenerateSummary(req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateActivityCopy(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.bridge.GenerateActivityCopy(req.Title, req.Description)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateSellingPoints(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.bridge.GenerateSellingPoints(req.Name, req.Description)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateActivityReport(c *gin.Context) {
	var req struct {
		ActivityName string                 `json:"activity_name" binding:"required"`
		Stats        map[string]interface{} `json:"stats"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.bridge.GenerateActivityReport(req.ActivityName, req.Stats)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, result)
}
