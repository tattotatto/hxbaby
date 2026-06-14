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

// BaaSActivityHandler serves activity/event API endpoints for mini-programs
// authenticated via X-API-Key.
type BaaSActivityHandler struct {
	db *gorm.DB
}

// NewBaaSActivityHandler creates a new BaaSActivityHandler.
func NewBaaSActivityHandler(db *gorm.DB) *BaaSActivityHandler {
	return &BaaSActivityHandler{db: db}
}

// ListActivities returns a paginated list of activities for the authenticated
// project, with optional status filter (draft/published/ended).
func (h *BaaSActivityHandler) ListActivities(c *gin.Context) {
	projectID := c.GetUint("project_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var activities []model.BActivity
	var total int64
	query := h.db.Where("project_id = ?", projectID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Model(&model.BActivity{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&activities)

	response.OK(c, gin.H{"items": activities, "total": total, "page": page, "page_size": pageSize})
}

// GetActivity returns a single activity by ID, scoped to the authenticated project.
func (h *BaaSActivityHandler) GetActivity(c *gin.Context) {
	projectID := c.GetUint("project_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的活动ID")
		return
	}

	var activity model.BActivity
	if err := h.db.Where("id = ? AND project_id = ?", id, projectID).First(&activity).Error; err != nil {
		response.Error(c, http.StatusNotFound, "活动不存在")
		return
	}
	response.OK(c, activity)
}

// SignupActivity creates a signup record for an activity, scoped to the
// authenticated project, and increments the activity's current_count.
func (h *BaaSActivityHandler) SignupActivity(c *gin.Context) {
	projectID := c.GetUint("project_id")
	activityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的活动ID")
		return
	}

	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Phone  string `json:"phone" binding:"required"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify activity exists and belongs to this project
	var activity model.BActivity
	if err := h.db.Where("id = ? AND project_id = ?", activityID, projectID).First(&activity).Error; err != nil {
		response.Error(c, http.StatusNotFound, "活动不存在")
		return
	}

	signup := &model.BActivitySignup{
		ActivityID: uint(activityID),
		ProjectID:  projectID,
		UserID:     req.UserID,
		Name:       req.Name,
		Phone:      req.Phone,
		Remark:     req.Remark,
	}

	if err := h.db.Create(signup).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "报名失败")
		return
	}

	// Increment current_count atomically
	h.db.Model(&activity).UpdateColumn("current_count", gorm.Expr("current_count + 1"))

	response.OK(c, signup)
}

// CheckinActivity marks a signup record as checked-in for the authenticated project.
func (h *BaaSActivityHandler) CheckinActivity(c *gin.Context) {
	projectID := c.GetUint("project_id")
	activityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的活动ID")
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var signup model.BActivitySignup
	if err := h.db.Where("activity_id = ? AND project_id = ? AND user_id = ?", activityID, projectID, req.UserID).
		First(&signup).Error; err != nil {
		response.Error(c, http.StatusNotFound, "报名记录不存在")
		return
	}

	now := time.Now()
	signup.CheckedIn = true
	signup.CheckedInAt = &now
	if err := h.db.Save(&signup).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "签到失败")
		return
	}

	response.OK(c, signup)
}
