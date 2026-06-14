package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/pkg/response"
)

type ChildHandler struct {
	childRepo *repository.ChildRepo
}

func NewChildHandler(repo *repository.ChildRepo) *ChildHandler {
	return &ChildHandler{childRepo: repo}
}

type CreateChildReq struct {
	Name      string `json:"name" binding:"required"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date" binding:"required"`
}

func (h *ChildHandler) Create(c *gin.Context) {
	var req CreateChildReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请填写完整信息")
		return
	}
	parentID := c.GetUint("user_id")
	tenantID := c.GetUint("tenant_id")

	birthDate, _ := time.Parse("2006-01-02", req.BirthDate)
	child := &model.Child{
		ParentID:  parentID,
		TenantID:  tenantID,
		Name:      req.Name,
		Gender:    req.Gender,
		BirthDate: birthDate,
	}
	if err := h.childRepo.Create(child); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}
	response.OK(c, child)
}

func (h *ChildHandler) List(c *gin.Context) {
	parentID := c.GetUint("user_id")
	children, err := h.childRepo.FindByParentID(parentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, children)
}

func (h *ChildHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	child, err := h.childRepo.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "儿童档案不存在")
		return
	}
	response.OK(c, child)
}

type AddGrowthReq struct {
	Date       string  `json:"date" binding:"required"`
	Height     float64 `json:"height"`
	Weight     float64 `json:"weight"`
	HeadCircum float64 `json:"head_circum"`
	Note       string  `json:"note"`
}

func (h *ChildHandler) AddGrowth(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	child, err := h.childRepo.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "儿童档案不存在")
		return
	}

	var req AddGrowthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请填写完整数据")
		return
	}

	record := model.GrowthRecord{
		Date: req.Date, Height: req.Height,
		Weight: req.Weight, HeadCircum: req.HeadCircum, Note: req.Note,
	}
	child.GrowthRecords = append(child.GrowthRecords, record)
	if err := h.childRepo.Update(child); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存失败")
		return
	}
	response.OK(c, child)
}
