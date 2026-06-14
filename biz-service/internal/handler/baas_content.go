package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/pkg/response"
)

// BaaSContentHandler serves content API endpoints for mini-programs
// authenticated via X-API-Key.
type BaaSContentHandler struct {
	repo *repository.CmsArticleRepo
}

// NewBaaSContentHandler creates a new BaaSContentHandler.
func NewBaaSContentHandler(db *gorm.DB) *BaaSContentHandler {
	return &BaaSContentHandler{repo: repository.NewCmsArticleRepo(db)}
}

// ListArticles returns a paginated list of published articles for the
// authenticated project, with optional category filtering.
func (h *BaaSContentHandler) ListArticles(c *gin.Context) {
	projectID := c.GetUint("project_id")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	articles, total, err := h.repo.List(projectID, category, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文章列表失败")
		return
	}
	response.OK(c, gin.H{
		"items":     articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetArticle returns a single article by ID, scoped to the authenticated
// project, and increments the view count.
func (h *BaaSContentHandler) GetArticle(c *gin.Context) {
	projectID := c.GetUint("project_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文章ID")
		return
	}

	article, err := h.repo.FindByID(uint(id), projectID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "文章不存在")
		return
	}
	// Increment view count (non-blocking)
	_ = h.repo.IncrementViewCount(uint(id))
	article.ViewCount++

	response.OK(c, article)
}
