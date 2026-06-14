package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/pkg/response"
)

// BaaSUserHandler serves user API endpoints for mini-programs
// authenticated via X-API-Key.
type BaaSUserHandler struct {
	db *gorm.DB
}

// NewBaaSUserHandler creates a new BaaSUserHandler.
func NewBaaSUserHandler(db *gorm.DB) *BaaSUserHandler {
	return &BaaSUserHandler{db: db}
}

// WxLogin is a stub for WeChat mini-program login. It accepts a code and
// returns a mock user + token. In production this would call the WeChat
// code2session API to resolve OpenID/UnionID.
func (h *BaaSUserHandler) WxLogin(c *gin.Context) {
	projectID := c.GetUint("project_id")

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Stub: in production, call WeChat code2session to get OpenID/UnionID,
	// then find-or-create the BUser. For now return a mock user.
	user := &model.BUser{
		ProjectID: projectID,
		OpenID:    fmt.Sprintf("mock_openid_%s", req.Code[:min(len(req.Code), 8)]),
		Nickname:  "微信用户",
	}
	if err := h.db.Where("open_id = ? AND project_id = ?", user.OpenID, projectID).
		FirstOrCreate(user).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "登录失败")
		return
	}

	// TODO: generate real JWT; for now return a placeholder token
	token := fmt.Sprintf("baas_token_%d_%d", user.ID, time.Now().Unix())

	response.OK(c, gin.H{
		"user":  user,
		"token": token,
	})
}

// GetProfile returns the current user's profile by user ID (from query param).
func (h *BaaSUserHandler) GetProfile(c *gin.Context) {
	projectID := c.GetUint("project_id")
	userID := c.Query("user_id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "缺少user_id参数")
		return
	}

	var user model.BUser
	if err := h.db.Where("id = ? AND project_id = ?", userID, projectID).First(&user).Error; err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	response.OK(c, user)
}

// UpdateProfile updates the current user's profile (nickname, avatar, phone).
func (h *BaaSUserHandler) UpdateProfile(c *gin.Context) {
	projectID := c.GetUint("project_id")

	var req struct {
		UserID   uint   `json:"user_id" binding:"required"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user model.BUser
	if err := h.db.Where("id = ? AND project_id = ?", req.UserID, projectID).First(&user).Error; err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if len(updates) > 0 {
		if err := h.db.Model(&user).Updates(updates).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "更新失败")
			return
		}
	}

	// Reload to return fresh data
	h.db.Where("id = ? AND project_id = ?", req.UserID, projectID).First(&user)
	response.OK(c, user)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

