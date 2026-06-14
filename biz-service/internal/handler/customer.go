package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/service"
	"github.com/hxbaby/biz-service/pkg/response"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

type CustomerRegisterReq struct {
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

func (h *CustomerHandler) Register(c *gin.Context) {
	var req CustomerRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	customer, token, err := h.svc.Register(req.Phone, req.Password, req.Name)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}
	response.OK(c, gin.H{"customer": customer, "token": token})
}

type CustomerLoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *CustomerHandler) Login(c *gin.Context) {
	var req CustomerLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	customer, token, err := h.svc.Login(req.Phone, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, gin.H{"customer": customer, "token": token})
}
