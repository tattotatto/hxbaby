package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/pkg/response"
)

// BaaSBookingHandler serves booking/reservation API endpoints for mini-programs
// authenticated via X-API-Key.
type BaaSBookingHandler struct {
	db *gorm.DB
}

// NewBaaSBookingHandler creates a new BaaSBookingHandler.
func NewBaaSBookingHandler(db *gorm.DB) *BaaSBookingHandler {
	return &BaaSBookingHandler{db: db}
}

// ListSlots returns available booking time slots for a given date (stub
// implementation — in production this would compute availability).
func (h *BaaSBookingHandler) ListSlots(c *gin.Context) {
	projectID := c.GetUint("project_id")
	date := c.Query("date") // expected format: 2026-06-14

	_ = projectID
	_ = date

	// Stub: return fixed time slots
	slots := []gin.H{
		{"time": "09:00", "available": true},
		{"time": "10:00", "available": true},
		{"time": "11:00", "available": false},
		{"time": "14:00", "available": true},
		{"time": "15:00", "available": true},
		{"time": "16:00", "available": true},
	}

	response.OK(c, gin.H{"date": date, "slots": slots})
}

// Create creates a new booking for the authenticated project.
func (h *BaaSBookingHandler) Create(c *gin.Context) {
	projectID := c.GetUint("project_id")

	var req struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Title   string `json:"title" binding:"required"`
		SlotStr string `json:"slot_time" binding:"required"` // ISO 8601 string
		Duration int   `json:"duration"`
		Name    string `json:"name" binding:"required"`
		Phone   string `json:"phone" binding:"required"`
		Remark  string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	slotTime, err := time.Parse(time.RFC3339, req.SlotStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的时间格式")
		return
	}

	duration := req.Duration
	if duration <= 0 {
		duration = 30
	}

	booking := &model.BBooking{
		ProjectID: projectID,
		UserID:    req.UserID,
		Title:     req.Title,
		SlotTime:  slotTime,
		Duration:  duration,
		Name:      req.Name,
		Phone:     req.Phone,
		Remark:    req.Remark,
		Status:    "confirmed",
	}

	if err := h.db.Create(booking).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "创建预约失败")
		return
	}
	response.OK(c, booking)
}

// List returns a paginated list of bookings for the authenticated project,
// with optional user_id and status filters.
func (h *BaaSBookingHandler) List(c *gin.Context) {
	projectID := c.GetUint("project_id")
	userID := c.Query("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var bookings []model.BBooking
	var total int64
	query := h.db.Where("project_id = ?", projectID)
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Model(&model.BBooking{}).Count(&total)
	query.Order("slot_time ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&bookings)

	response.OK(c, gin.H{"items": bookings, "total": total, "page": page, "page_size": pageSize})
}
