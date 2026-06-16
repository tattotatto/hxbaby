package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/pkg/response"
)

type KnowledgeHandler struct {
	aiServiceURL string
	httpClient   *http.Client
}

func NewKnowledgeHandler(aiServiceURL string) *KnowledgeHandler {
	return &KnowledgeHandler{
		aiServiceURL: strings.TrimRight(aiServiceURL, "/"),
		httpClient:   &http.Client{},
	}
}

// proxyToAI 转发请求到 Python AI 服务
func (h *KnowledgeHandler) proxyToAI(c *gin.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, h.aiServiceURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create proxy request failed: %w", err)
	}
	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	return h.httpClient.Do(req)
}

// Upload 上传文档入库 — 直接转发 multipart
func (h *KnowledgeHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "无法读取文件")
		return
	}
	defer f.Close()

	// 用管道转发 multipart
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		defer f.Close()
		fmt.Fprintf(pw, "--boundary\r\n")
		fmt.Fprintf(pw, "Content-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n", file.Filename)
		fmt.Fprintf(pw, "Content-Type: application/octet-stream\r\n\r\n")
		io.Copy(pw, f)
		fmt.Fprintf(pw, "\r\n--boundary--\r\n")
	}()

	req, _ := http.NewRequest("POST", h.aiServiceURL+"/ai/knowledge/upload", pr)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")
	resp, err := h.httpClient.Do(req)
	if err != nil {
		response.Error(c, http.StatusBadGateway, fmt.Sprintf("AI 服务不可用: %v", err))
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// ListDocuments 获取知识库文档列表
func (h *KnowledgeHandler) ListDocuments(c *gin.Context) {
	resp, err := h.proxyToAI(c, "GET", "/ai/knowledge/documents", nil)
	if err != nil {
		response.Error(c, http.StatusBadGateway, fmt.Sprintf("AI 服务不可用: %v", err))
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// DeleteDocument 删除文档
func (h *KnowledgeHandler) DeleteDocument(c *gin.Context) {
	source := c.Param("source")
	encoded := url.PathEscape(source)
	resp, err := h.proxyToAI(c, "DELETE", "/ai/knowledge/documents/"+encoded, nil)
	if err != nil {
		response.Error(c, http.StatusBadGateway, fmt.Sprintf("AI 服务不可用: %v", err))
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

// GetStats 获取知识库统计
func (h *KnowledgeHandler) GetStats(c *gin.Context) {
	resp, err := h.proxyToAI(c, "GET", "/ai/knowledge/stats", nil)
	if err != nil {
		response.Error(c, http.StatusBadGateway, fmt.Sprintf("AI 服务不可用: %v", err))
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

