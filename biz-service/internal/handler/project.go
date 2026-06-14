package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/service"
	"github.com/hxbaby/biz-service/pkg/response"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req service.CreateProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	customerID := c.GetUint("user_id")
	project, err := h.svc.Create(customerID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	customerID := c.GetUint("user_id")
	projects, err := h.svc.List(customerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取项目列表失败")
		return
	}
	response.OK(c, projects)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	project, err := h.svc.Get(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "项目不存在")
		return
	}
	response.OK(c, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		Modules     []string               `json:"modules"`
		BrandConfig map[string]interface{} `json:"brand_config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.svc.UpdateModules(uint(id), req.Modules, req.BrandConfig); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, nil)
}
