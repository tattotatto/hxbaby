package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/service"
	"github.com/hxbaby/biz-service/pkg/response"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

type LoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请输入手机号和密码")
		return
	}
	user, token, err := h.svc.Login(req.Phone, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, gin.H{"user": user, "token": token})
}

type RegisterReq struct {
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请填写完整信息")
		return
	}
	// 注册时 tenantID 默认为1（后续从邀请链接获取）
	user, token, err := h.svc.Register(req.Phone, req.Password, req.Name, 1)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}
	response.OK(c, gin.H{"user": user, "token": token})
}
