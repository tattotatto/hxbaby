package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/pkg/response"
)

// BaaSShopHandler serves shop/ecommerce API endpoints for mini-programs
// authenticated via X-API-Key.
type BaaSShopHandler struct {
	db *gorm.DB
}

// NewBaaSShopHandler creates a new BaaSShopHandler.
func NewBaaSShopHandler(db *gorm.DB) *BaaSShopHandler {
	return &BaaSShopHandler{db: db}
}

// ListProducts returns a paginated list of published products for the
// authenticated project, with optional category filtering.
func (h *BaaSShopHandler) ListProducts(c *gin.Context) {
	projectID := c.GetUint("project_id")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var products []model.BProduct
	var total int64
	query := h.db.Where("project_id = ? AND status = ?", projectID, "on")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query.Model(&model.BProduct{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

	response.OK(c, gin.H{"items": products, "total": total, "page": page, "page_size": pageSize})
}

// GetProduct returns a single product by ID, scoped to the authenticated project.
func (h *BaaSShopHandler) GetProduct(c *gin.Context) {
	projectID := c.GetUint("project_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var product model.BProduct
	if err := h.db.Where("id = ? AND project_id = ?", id, projectID).First(&product).Error; err != nil {
		response.Error(c, http.StatusNotFound, "商品不存在")
		return
	}
	response.OK(c, product)
}

// CreateOrder creates a new order for the authenticated project.
func (h *BaaSShopHandler) CreateOrder(c *gin.Context) {
	projectID := c.GetUint("project_id")
	// TODO: get userID from BaaS user token (placeholder: from request body)
	var req struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Items   string `json:"items" binding:"required"` // JSON string array
		Address string `json:"address"`
		Remark  string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Generate simple order number
	orderNo := fmt.Sprintf("BO%d%d", time.Now().UnixMilli(), projectID)

	order := &model.BOrder{
		ProjectID:   projectID,
		UserID:      req.UserID,
		OrderNo:     orderNo,
		Items:       req.Items,
		TotalAmount: 0, // In production, calculate from items
		Status:      "pending",
		Address:     req.Address,
		Remark:      req.Remark,
	}
	if err := h.db.Create(order).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建订单失败")
		return
	}
	response.OK(c, order)
}

// ListOrders returns a paginated list of orders for the authenticated project,
// with optional user_id and status filters.
func (h *BaaSShopHandler) ListOrders(c *gin.Context) {
	projectID := c.GetUint("project_id")
	userID := c.Query("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var orders []model.BOrder
	var total int64
	query := h.db.Where("project_id = ?", projectID)
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Model(&model.BOrder{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

	response.OK(c, gin.H{"items": orders, "total": total, "page": page, "page_size": pageSize})
}

// GetOrder returns a single order by ID, scoped to the authenticated project.
func (h *BaaSShopHandler) GetOrder(c *gin.Context) {
	projectID := c.GetUint("project_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var order model.BOrder
	if err := h.db.Where("id = ? AND project_id = ?", id, projectID).First(&order).Error; err != nil {
		response.Error(c, http.StatusNotFound, "订单不存在")
		return
	}
	response.OK(c, order)
}
