package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/pkg/response"
)

type ProductHandler struct {
	repo *repository.ProductRepo
}

func NewProductHandler(repo *repository.ProductRepo) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) List(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	products, err := h.repo.FindByTenantID(tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, products)
}

func (h *ProductHandler) Match(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	symptoms := c.Query("symptoms")
	if symptoms == "" {
		response.Error(c, http.StatusBadRequest, "请提供症状关键词")
		return
	}
	products, err := h.repo.SearchBySymptoms(tenantID, strings.Split(symptoms, ","))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "匹配失败")
		return
	}
	response.OK(c, products)
}
