package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/pkg/response"
)

type BuildHandler struct {
	buildRepo   *repository.BuildTaskRepo
	projectRepo *repository.ProjectRepo
	codegenURL  string
}

func NewBuildHandler(buildRepo *repository.BuildTaskRepo, projectRepo *repository.ProjectRepo, codegenURL string) *BuildHandler {
	return &BuildHandler{
		buildRepo:   buildRepo,
		projectRepo: projectRepo,
		codegenURL:  codegenURL,
	}
}

// TriggerBuild starts a code generation build for a project
func (h *BuildHandler) TriggerBuild(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}
	customerID := c.GetUint("user_id")

	// Load project
	project, err := h.projectRepo.FindByID(uint(id))
	if err != nil || project.CustomerID != customerID {
		response.Error(c, http.StatusNotFound, "项目不存在")
		return
	}

	// Create BuildTask
	task := &model.BuildTask{
		ProjectID:   project.ID,
		TriggeredBy: customerID,
		Status:      "pending",
	}
	if err := h.buildRepo.Create(task); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建构建任务失败")
		return
	}

	// Call codegen-service asynchronously (goroutine)
	go h.executeBuild(task, project)

	response.OK(c, gin.H{
		"build_id": task.ID,
		"status":   task.Status,
	})
}

// executeBuild calls the codegen-service and updates the build task
func (h *BuildHandler) executeBuild(task *model.BuildTask, project *model.MiniappProject) {
	startTime := time.Now()

	// Parse modules JSON string to []string
	var modules []string
	json.Unmarshal([]byte(project.Modules), &modules)

	// Parse brand config JSON string to map
	var brandConfig map[string]interface{}
	json.Unmarshal([]byte(project.BrandConfig), &brandConfig)

	// Build request payload for codegen-service
	payload := map[string]interface{}{
		"project": map[string]interface{}{
			"id":           project.ID,
			"name":         project.Name,
			"modules":      modules,
			"brand_config": brandConfig,
			"api_key":      project.APIKey,
			"wx_app_id":    project.WxAppID,
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(
		h.codegenURL+"/api/build",
		"application/json",
		bytes.NewReader(body),
	)

	// Update task status
	task.DurationMs = time.Since(startTime).Milliseconds()

	if err != nil {
		task.Status = "failed"
		task.ErrorLog = fmt.Sprintf("代码生成服务调用失败: %v", err)
		h.buildRepo.Update(task)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		task.Status = "failed"
		errBody, _ := io.ReadAll(resp.Body)
		task.ErrorLog = fmt.Sprintf("代码生成失败 [%d]: %s", resp.StatusCode, string(errBody))
		h.buildRepo.Update(task)
		return
	}

	var result struct {
		TaskID    string   `json:"task_id"`
		Status    string   `json:"status"`
		ZipPath   string   `json:"zip_path"`
		MD5       string   `json:"md5"`
		SizeBytes int64    `json:"size_bytes"`
		Warnings  []string `json:"warnings"`
		Error     string   `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		task.Status = "failed"
		task.ErrorLog = fmt.Sprintf("解析代码生成响应失败: %v", err)
		h.buildRepo.Update(task)
		return
	}

	if result.Status == "done" {
		now := time.Now()
		task.Status = "done"
		task.OutputZipURL = result.ZipPath
		task.OutputMD5 = result.MD5
		task.CompletedAt = &now
	} else {
		task.Status = "failed"
		task.ErrorLog = result.Error
	}
	h.buildRepo.Update(task)
}

// GetBuildStatus returns the status of a build task
func (h *BuildHandler) GetBuildStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的构建任务ID")
		return
	}

	task, err := h.buildRepo.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "构建任务不存在")
		return
	}

	response.OK(c, task)
}

// DownloadBuild returns the ZIP file for download
func (h *BuildHandler) DownloadBuild(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的构建任务ID")
		return
	}

	task, err := h.buildRepo.FindByID(uint(id))
	if err != nil || task.Status != "done" {
		response.Error(c, http.StatusNotFound, "构建未完成或不存在")
		return
	}

	// For now, serve the file directly
	// In production, this would redirect to OSS signed URL
	if task.OutputZipURL != "" {
		c.File(task.OutputZipURL)
		return
	}

	response.Error(c, http.StatusNotFound, "下载文件不存在")
}

// GetBuildHistory returns all build tasks for a project
func (h *BuildHandler) GetBuildHistory(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	tasks, err := h.buildRepo.FindByProjectID(uint(projectID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取构建历史失败")
		return
	}

	response.OK(c, tasks)
}
