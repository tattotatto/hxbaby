# 多用户小程序工厂平台 — 详细实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建SaaS小程序工厂平台MVP，客户注册后可自由组合8大功能模块、一键生成UniApp源码包、通过全托管BaaS获得AI增强的后端服务。

**Architecture:** Go扩展(工厂管理+BaaS API) + Node.js微服务(代码生成引擎) + React管理前端(Ant Design 5)，复用AI系统的PostgreSQL/Redis/OSS/Python AI。

**Tech Stack:** Go 1.22+ (Gin扩展), Node.js 20+ (Express + Handlebars + archiver), React 18 + Ant Design 5, UniApp (Vue 3)

---

## File Structure (新增文件)

```
hxbaby/
│
├── factory-admin/                  # React 管理前端 (新增)
│   ├── package.json
│   ├── vite.config.ts
│   ├── index.html
│   ├── src/
│   │   ├── main.tsx
│   │   ├── App.tsx
│   │   ├── api/                   # API 调用封装
│   │   │   ├── client.ts          # axios 实例
│   │   │   ├── auth.ts
│   │   │   ├── project.ts
│   │   │   ├── build.ts
│   │   │   └── ai.ts              # AI增强API
│   │   ├── pages/
│   │   │   ├── Login.tsx
│   │   │   ├── Register.tsx
│   │   │   ├── Dashboard.tsx      # 项目列表
│   │   │   ├── ProjectCreate.tsx  # 创建向导
│   │   │   ├── ProjectConfig.tsx  # 模块选择+品牌配置
│   │   │   ├── ProjectBuild.tsx   # 构建&下载
│   │   │   ├── Guide.tsx          # 接入指引
│   │   │   └── Billing.tsx        # 套餐计费
│   │   ├── components/
│   │   │   ├── ModuleSelector.tsx # 模块勾选组件
│   │   │   ├── BrandConfig.tsx    # 品牌配置表单
│   │   │   ├── BuildProgress.tsx  # 构建进度
│   │   │   └── Layout.tsx         # 管理后台布局
│   │   └── hooks/
│   │       └── useAuth.ts
│   └── Dockerfile
│
├── codegen-service/               # Node.js 代码生成引擎 (新增)
│   ├── package.json
│   ├── server.js                  # Express 入口
│   ├── scripts/
│   │   ├── compose.js             # 编排主入口
│   │   ├── resolver.js            # 依赖解析+拓扑排序
│   │   ├── injector.js            # 配置注入
│   │   ├── packager.js            # ZIP打包+OSS上传
│   │   └── validator.js           # 合法性校验
│   ├── templates/
│   │   ├── base/                  # 基础骨架模板
│   │   │   ├── App.vue.hbs
│   │   │   ├── main.js.hbs
│   │   │   ├── manifest.json.hbs
│   │   │   ├── pages.json.hbs
│   │   │   ├── config.js.hbs
│   │   │   ├── common/
│   │   │   │   ├── baas-sdk.js.hbs       # BaaS SDK 模板
│   │   │   │   └── request.js.hbs        # HTTP 请求封装
│   │   │   └── components/
│   │   │       └── tabbar.vue.hbs
│   │   └── modules/
│   │       ├── cms/               # 📝 内容模块
│   │       │   ├── module.json
│   │       │   ├── pages/article-list.vue.hbs
│   │       │   ├── pages/article-detail.vue.hbs
│   │       │   └── api/content.js.hbs
│   │       ├── ai-advisor/        # 🤖 AI顾问模块
│   │       │   ├── module.json
│   │       │   ├── pages/chat.vue.hbs
│   │       │   ├── pages/symptom-check.vue.hbs
│   │       │   ├── pages/growth-record.vue.hbs
│   │       │   └── api/ai.js.hbs
│   │       ├── shop/              # 🛒 电商模块
│   │       │   ├── module.json
│   │       │   ├── pages/product-list.vue.hbs
│   │       │   ├── pages/product-detail.vue.hbs
│   │       │   ├── pages/cart.vue.hbs
│   │       │   ├── pages/checkout.vue.hbs
│   │       │   ├── pages/order-list.vue.hbs
│   │       │   └── api/shop.js.hbs
│   │       ├── activity/          # 🎪 活动模块(AI增强)
│   │       │   ├── module.json
│   │       │   ├── pages/activity-list.vue.hbs
│   │       │   ├── pages/activity-detail.vue.hbs
│   │       │   ├── pages/activity-signup.vue.hbs
│   │       │   ├── components/ai-copy-gen.vue.hbs
│   │       │   ├── components/ai-poster.vue.hbs
│   │       │   └── api/activity.js.hbs
│   │       ├── booking/           # 📅 预约模块
│   │       ├── member/            # 👑 会员模块
│   │       └── analytics/         # 📊 数据模块
│   ├── Dockerfile
│   └── tests/
│       ├── resolver.test.js
│       ├── injector.test.js
│       └── compose.test.js
│
├── biz-service/                   # Go 业务服务 (扩展)
│   ├── internal/
│   │   ├── model/
│   │   │   ├── customer.go        # 新增
│   │   │   ├── miniapp_project.go # 新增
│   │   │   ├── build_task.go      # 新增
│   │   │   ├── cms_article.go     # 新增
│   │   │   ├── b_activity.go      # 新增
│   │   │   ├── b_product.go       # 新增
│   │   │   ├── b_order.go         # 新增
│   │   │   ├── b_user.go          # 新增
│   │   │   └── b_booking.go       # 新增
│   │   ├── handler/
│   │   │   ├── customer.go        # 新增：工厂管理
│   │   │   ├── project.go         # 新增：项目管理
│   │   │   ├── baas_content.go    # 新增：BaaS内容API
│   │   │   ├── baas_shop.go       # 新增：BaaS电商API
│   │   │   ├── baas_activity.go   # 新增：BaaS活动API
│   │   │   ├── baas_user.go       # 新增：BaaS用户API
│   │   │   ├── baas_ai.go         # 新增：BaaS AI桥接API
│   │   │   └── ai_enhance.go      # 新增：AI增强API
│   │   ├── service/
│   │   │   ├── customer.go
│   │   │   ├── project.go
│   │   │   ├── baas_cms.go
│   │   │   ├── baas_shop.go
│   │   │   ├── baas_activity.go
│   │   │   ├── ai_bridge.go       # AI Bridge 统一调用
│   │   │   └── wechat_pay.go      # 微信支付
│   │   ├── repository/
│   │   │   ├── customer.go
│   │   │   ├── project.go
│   │   │   ├── build_task.go
│   │   │   ├── cms_article.go
│   │   │   ├── b_activity.go
│   │   │   ├── b_product.go
│   │   │   ├── b_order.go
│   │   │   └── b_user.go
│   │   └── middleware/
│   │       ├── baas_auth.go       # 新增：BaaS API Key鉴权
│   │       └── project_isolate.go # 新增：项目数据隔离
│   └── pkg/
│       └── apikey/
│           └── apikey.go          # API Key生成/校验
```

---

## Phase 1: Go扩展 + React管理后台 (Week 1-2)

### Task 1: Go工厂管理数据模型与迁移

**Files:**
- Create: `biz-service/internal/model/customer.go`
- Create: `biz-service/internal/model/miniapp_project.go`
- Create: `biz-service/internal/model/build_task.go`

- [ ] **Step 1: 创建 customer.go**

```go
package model

import "time"

type Customer struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"size:100;not null" json:"name"`
    Phone       string    `gorm:"uniqueIndex;size:20;not null" json:"phone"`
    Email       string    `gorm:"size:200" json:"email,omitempty"`
    PasswordHash string   `gorm:"size:255;not null" json:"-"`
    CompanyName string    `gorm:"size:200" json:"company_name,omitempty"`
    WxUnionID   string    `gorm:"size:100" json:"wx_unionid,omitempty"`
    Plan        string    `gorm:"size:20;default:'free'" json:"plan"` // free/basic/pro/enterprise
    PlanExpiresAt *time.Time `json:"plan_expires_at,omitempty"`
    MaxProjects int       `gorm:"default:1" json:"max_projects"`
    Status      string    `gorm:"size:20;default:'active'" json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

- [ ] **Step 2: 创建 miniapp_project.go**

```go
package model

import (
    "database/sql/driver"
    "encoding/json"
    "time"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) { return json.Marshal(j) }
func (j *JSONB) Scan(value interface{}) error {
    if value == nil { return nil }
    return json.Unmarshal(value.([]byte), j)
}

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) { return json.Marshal(s) }
func (s *StringSlice) Scan(value interface{}) error {
    if value == nil { return nil }
    return json.Unmarshal(value.([]byte), s)
}

type MiniappProject struct {
    ID          uint        `gorm:"primaryKey" json:"id"`
    CustomerID  uint        `gorm:"index;not null" json:"customer_id"`
    Name        string      `gorm:"size:100;not null" json:"name"`
    Description string      `gorm:"type:text" json:"description,omitempty"`
    WxAppID     string      `gorm:"size:50" json:"wx_app_id,omitempty"`
    WxAppSecret string      `gorm:"size:255" json:"-"` // 加密存储
    WxMchID     string      `gorm:"size:50" json:"wx_mch_id,omitempty"`
    Modules     StringSlice `gorm:"type:jsonb;default:'[]'" json:"modules"`
    BrandConfig JSONB       `gorm:"type:jsonb;default:'{}'" json:"brand_config"`
    APIKey      string      `gorm:"uniqueIndex;size:64;not null" json:"api_key"`
    APISecret   string      `gorm:"size:64;not null" json:"-"`
    Status      string      `gorm:"size:20;default:'draft'" json:"status"`
    Domain      string      `gorm:"size:200" json:"domain,omitempty"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}
```

- [ ] **Step 3: 创建 build_task.go**

```go
package model

import "time"

type BuildTask struct {
    ID               uint      `gorm:"primaryKey" json:"id"`
    ProjectID        uint      `gorm:"index;not null" json:"project_id"`
    TriggeredBy      uint      `json:"triggered_by"`
    ModulesSnapshot  JSONB     `gorm:"type:jsonb" json:"modules_snapshot"`
    ConfigSnapshot   JSONB     `gorm:"type:jsonb" json:"config_snapshot"`
    Status           string    `gorm:"size:20;default:'pending'" json:"status"`
    // pending/resolving/composing/injecting/packaging/done/failed
    OutputZipURL     string    `gorm:"size:500" json:"output_zip_url,omitempty"`
    OutputMD5        string    `gorm:"size:64" json:"output_md5,omitempty"`
    ErrorLog         string    `gorm:"type:text" json:"error_log,omitempty"`
    DurationMs       int64     `json:"duration_ms,omitempty"`
    CreatedAt        time.Time `json:"created_at"`
    CompletedAt      *time.Time `json:"completed_at,omitempty"`
}
```

- [ ] **Step 4: 更新 AutoMigrate**

在 `internal/config/database.go` 的 AutoMigrate 函数中添加：
```go
DB.AutoMigrate(
    // ... existing models
    &model.Customer{},
    &model.MiniappProject{},
    &model.BuildTask{},
    &model.CmsArticle{},
    &model.BActivity{},
    &model.BProduct{},
    &model.BOrder{},
    &model.BUser{},
    &model.BBooking{},
)
```

- [ ] **Step 5: Commit**

```bash
git add biz-service/internal/model/customer.go biz-service/internal/model/miniapp_project.go biz-service/internal/model/build_task.go
git commit -m "feat: add factory management data models (customer, project, build_task)"
```

---

### Task 2: 客户注册/登录 API

**Files:**
- Create: `biz-service/internal/repository/customer.go`
- Create: `biz-service/internal/service/customer.go`
- Create: `biz-service/internal/handler/customer.go`

- [ ] **Step 1: 写出 repository/customer.go**

```go
package repository

import (
    "github.com/hxbaby/biz-service/internal/model"
    "gorm.io/gorm"
)

type CustomerRepo struct {
    db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
    return &CustomerRepo{db: db}
}

func (r *CustomerRepo) Create(c *model.Customer) error {
    return r.db.Create(c).Error
}

func (r *CustomerRepo) FindByPhone(phone string) (*model.Customer, error) {
    var c model.Customer
    err := r.db.Where("phone = ?", phone).First(&c).Error
    if err != nil {
        return nil, err
    }
    return &c, nil
}

func (r *CustomerRepo) FindByID(id uint) (*model.Customer, error) {
    var c model.Customer
    err := r.db.First(&c, id).Error
    if err != nil {
        return nil, err
    }
    return &c, nil
}

func (r *CustomerRepo) Update(c *model.Customer) error {
    return r.db.Save(c).Error
}
```

- [ ] **Step 2: 写出 service/customer.go**

```go
package service

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/hxbaby/biz-service/internal/model"
    "github.com/hxbaby/biz-service/internal/repository"
    "github.com/hxbaby/biz-service/pkg/jwt"
)

type CustomerService struct {
    repo      *repository.CustomerRepo
    jwtSecret string
}

func NewCustomerService(repo *repository.CustomerRepo, jwtSecret string) *CustomerService {
    return &CustomerService{repo: repo, jwtSecret: jwtSecret}
}

func (s *CustomerService) Register(phone, password, name string) (*model.Customer, string, error) {
    existing, _ := s.repo.FindByPhone(phone)
    if existing != nil {
        return nil, "", errors.New("手机号已注册")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, "", err
    }

    c := &model.Customer{
        Phone:       phone,
        PasswordHash: string(hash),
        Name:        name,
        Plan:        "free",
        MaxProjects: 1,
        Status:      "active",
    }
    if err := s.repo.Create(c); err != nil {
        return nil, "", err
    }

    token, err := jwt.GenerateToken(c.ID, 0, "customer", s.jwtSecret)
    return c, token, err
}

func (s *CustomerService) Login(phone, password string) (*model.Customer, string, error) {
    c, err := s.repo.FindByPhone(phone)
    if err != nil {
        return nil, "", errors.New("手机号未注册")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(password)); err != nil {
        return nil, "", errors.New("密码错误")
    }
    token, err := jwt.GenerateToken(c.ID, 0, "customer", s.jwtSecret)
    return c, token, err
}
```

- [ ] **Step 3: 写出 handler/customer.go**

```go
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

type RegisterReq struct {
    Phone    string `json:"phone" binding:"required,len=11"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required"`
}

func (h *CustomerHandler) Register(c *gin.Context) {
    var req RegisterReq
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

type LoginReq struct {
    Phone    string `json:"phone" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (h *CustomerHandler) Login(c *gin.Context) {
    var req LoginReq
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
```

- [ ] **Step 4: 注册路由**

```go
// router.go
customerH := handler.NewCustomerHandler(service.NewCustomerService(
    repository.NewCustomerRepo(config.DB),
    cfg.JWTSecret,
))
// 公开接口
r.POST("/api/v1/auth/register", customerH.Register)
r.POST("/api/v1/auth/login", customerH.Login)
```

- [ ] **Step 5: Commit**

---

### Task 3: 项目管理 API (CRUD + API Key生成)

**Files:**
- Create: `biz-service/pkg/apikey/apikey.go`
- Create: `biz-service/internal/repository/project.go`
- Create: `biz-service/internal/service/project.go`
- Create: `biz-service/internal/handler/project.go`

- [ ] **Step 1: 实现 API Key 生成工具 pkg/apikey/apikey.go**

```go
package apikey

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

func GenerateKey() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return "ak_" + hex.EncodeToString(bytes), nil
}

func GenerateSecret() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func HashSecret(secret string) string {
    h := sha256.Sum256([]byte(secret))
    return hex.EncodeToString(h[:])
}
```

- [ ] **Step 2: 写出 repository/project.go**

```go
package repository

import (
    "github.com/hxbaby/biz-service/internal/model"
    "gorm.io/gorm"
)

type ProjectRepo struct {
    db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
    return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(p *model.MiniappProject) error {
    return r.db.Create(p).Error
}

func (r *ProjectRepo) FindByCustomerID(customerID uint) ([]model.MiniappProject, error) {
    var projects []model.MiniappProject
    err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&projects).Error
    return projects, err
}

func (r *ProjectRepo) FindByID(id uint) (*model.MiniappProject, error) {
    var p model.MiniappProject
    err := r.db.First(&p, id).Error
    if err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *ProjectRepo) CountByCustomerID(customerID uint) (int64, error) {
    var count int64
    err := r.db.Model(&model.MiniappProject{}).Where("customer_id = ?", customerID).Count(&count).Error
    return count, err
}

func (r *ProjectRepo) Update(p *model.MiniappProject) error {
    return r.db.Save(p).Error
}
```

- [ ] **Step 3: 写出 service/project.go**

```go
package service

import (
    "errors"
    "github.com/hxbaby/biz-service/internal/model"
    "github.com/hxbaby/biz-service/internal/repository"
    "github.com/hxbaby/biz-service/pkg/apikey"
)

type ProjectService struct {
    projectRepo  *repository.ProjectRepo
    customerRepo *repository.CustomerRepo
}

func NewProjectService(pr *repository.ProjectRepo, cr *repository.CustomerRepo) *ProjectService {
    return &ProjectService{projectRepo: pr, customerRepo: cr}
}

type CreateProjectReq struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Modules     []string `json:"modules"`
}

func (s *ProjectService) Create(customerID uint, req CreateProjectReq) (*model.MiniappProject, error) {
    // 检查项目数限制
    customer, err := s.customerRepo.FindByID(customerID)
    if err != nil {
        return nil, errors.New("客户不存在")
    }
    count, _ := s.projectRepo.CountByCustomerID(customerID)
    if int(count) >= customer.MaxProjects {
        return nil, errors.New("已达到项目数量上限，请升级套餐")
    }

    apiKey, _ := apikey.GenerateKey()
    apiSecret, _ := apikey.GenerateSecret()

    // 确保至少包含基础模块
    modules := req.Modules
    hasBase := false
    for _, m := range modules {
        if m == "base" { hasBase = true; break }
    }
    if !hasBase {
        modules = append([]string{"base"}, modules...)
    }

    project := &model.MiniappProject{
        CustomerID:  customerID,
        Name:        req.Name,
        Description: req.Description,
        Modules:     modules,
        APIKey:      apiKey,
        APISecret:   apikey.HashSecret(apiSecret),
        BrandConfig: model.JSONB{
            "appName":      req.Name,
            "primaryColor": "#4caf50",
        },
        Status: "draft",
    }
    if err := s.projectRepo.Create(project); err != nil {
        return nil, err
    }
    // 返回时带上明文 secret（仅此一次）
    project.APISecret = apiSecret
    return project, nil
}

func (s *ProjectService) List(customerID uint) ([]model.MiniappProject, error) {
    return s.projectRepo.FindByCustomerID(customerID)
}

func (s *ProjectService) Get(id uint) (*model.MiniappProject, error) {
    return s.projectRepo.FindByID(id)
}

func (s *ProjectService) UpdateModules(id uint, modules []string, brandConfig map[string]interface{}) error {
    p, err := s.projectRepo.FindByID(id)
    if err != nil { return err }
    if modules != nil { p.Modules = modules }
    if brandConfig != nil { p.BrandConfig = model.JSONB(brandConfig) }
    return s.projectRepo.Update(p)
}
```

- [ ] **Step 4: 写出 handler/project.go**

```go
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
```

- [ ] **Step 5: Commit**

---

## Phase 2: Node.js 代码生成引擎 (Week 2-3)

### Task 4: 代码生成引擎骨架 + 依赖解析器

**Files:**
- Create: `codegen-service/package.json`
- Create: `codegen-service/server.js`
- Create: `codegen-service/scripts/resolver.js`
- Create: `codegen-service/scripts/validator.js`
- Create: `codegen-service/tests/resolver.test.js`

- [ ] **Step 1: 初始化项目**

```bash
mkdir -p codegen-service/{scripts,templates/{base,modules},tests}
cd codegen-service
npm init -y
npm install express handlebars archiver uuid
npm install --save-dev jest
```

- [ ] **Step 2: 创建 package.json scripts**

```json
{
  "name": "codegen-service",
  "version": "0.1.0",
  "main": "server.js",
  "scripts": {
    "start": "node server.js",
    "test": "jest --verbose"
  },
  "dependencies": {
    "archiver": "^6.0.0",
    "express": "^4.19.0",
    "handlebars": "^4.7.8",
    "uuid": "^9.0.0"
  },
  "devDependencies": {
    "jest": "^29.7.0"
  }
}
```

- [ ] **Step 3: 创建 resolver.js 测试 tests/resolver.test.js**

```javascript
const { resolveModules, detectConflicts, DEPENDENCY_GRAPH } = require('../scripts/resolver');

describe('Module Resolver', () => {
  test('resolves modules with auto-dependency addition', () => {
    // 选了shop但没选微信支付基础配置
    const selected = ['base', 'shop'];
    const resolved = resolveModules(selected);
    // shop 依赖 base, base已在
    expect(resolved).toContain('base');
    expect(resolved).toContain('shop');
  });

  test('auto-adds missing base module', () => {
    const selected = ['cms'];
    const resolved = resolveModules(selected);
    expect(resolved).toContain('base'); // 自动补齐
    expect(resolved).toContain('cms');
  });

  test('detects conflicts', () => {
    // 选了电商但没有选微信支付的配置要求 -> 仅警告不阻止
    const warnings = detectConflicts(['base', 'shop']);
    // shop 需要微信支付，缺失则产生警告
    const hasPaymentWarning = warnings.some(w => w.includes('支付'));
    expect(hasPaymentWarning).toBe(true);
  });

  test('topological sort ensures dependencies first', () => {
    // shop 依赖 base -> base 应在 shop 前面
    const order = resolveModules(['shop', 'analytics', 'base']);
    const baseIdx = order.indexOf('base');
    const shopIdx = order.indexOf('shop');
    expect(baseIdx).toBeLessThan(shopIdx);
  });

  test('all modules in DEPENDENCY_GRAPH have valid dependencies', () => {
    const validModules = Object.keys(DEPENDENCY_GRAPH);
    for (const [mod, deps] of Object.entries(DEPENDENCY_GRAPH)) {
      for (const dep of deps) {
        expect(validModules).toContain(dep);
      }
    }
  });
});
```

- [ ] **Step 4: 运行测试确认失败**

```bash
cd codegen-service
npm test
# 预期: FAIL (模块不存在)
```

- [ ] **Step 5: 实现 scripts/resolver.js**

```javascript
// 模块依赖图
const DEPENDENCY_GRAPH = {
  'base':       [],
  'cms':        ['base'],
  'ai-advisor': ['base'],
  'shop':       ['base'],
  'activity':   ['base'],
  'booking':    ['base'],
  'member':     ['base'],
  'analytics':  ['ai-advisor'],
};

// 模块冲突/警告规则
const WARNING_RULES = [
  {
    // 选了电商但没有支付能力：仅警告
    condition: (mods) => mods.includes('shop') && !mods.includes('shop'),
    message: '电商模块需要配置微信支付，请确保已在微信支付商户平台完成接入',
  },
];

/**
 * 解析模块列表，自动补齐缺失依赖，返回拓扑排序后的模块列表
 */
function resolveModules(selected) {
  if (!selected.includes('base')) {
    selected = ['base', ...selected];
  }

  // BFS 补齐所有依赖
  const resolved = new Set();
  const queue = [...selected];

  while (queue.length > 0) {
    const mod = queue.shift();
    if (resolved.has(mod)) continue;
    resolved.add(mod);

    const deps = DEPENDENCY_GRAPH[mod] || [];
    for (const dep of deps) {
      if (!resolved.has(dep)) {
        queue.push(dep);
      }
    }
  }

  // 拓扑排序
  return topologicalSort([...resolved]);
}

function topologicalSort(modules) {
  const sorted = [];
  const visited = new Set();
  const temp = new Set();

  function visit(mod) {
    if (temp.has(mod)) return; // 忽略循环依赖
    if (visited.has(mod)) return;
    temp.add(mod);
    const deps = DEPENDENCY_GRAPH[mod] || [];
    for (const dep of deps) {
      visit(dep);
    }
    temp.delete(mod);
    visited.add(mod);
    sorted.push(mod);
  }

  for (const mod of modules) {
    visit(mod);
  }

  return sorted;
}

function detectConflicts(modules) {
  const warnings = [];
  for (const rule of WARNING_RULES) {
    if (rule.condition(modules)) {
      warnings.push(rule.message);
    }
  }
  return warnings;
}

module.exports = { resolveModules, detectConflicts, DEPENDENCY_GRAPH };
```

- [ ] **Step 6: 运行测试确认通过**

```bash
npm test
# 预期: PASS
```

- [ ] **Step 7: Commit**

```bash
git add codegen-service/
git commit -m "feat: add codegen service skeleton with module resolver"
```

---

### Task 5: 模块声明规范 + 配置注入器

**Files:**
- Create: `codegen-service/templates/modules/cms/module.json`
- Create: `codegen-service/templates/modules/ai-advisor/module.json`
- Create: `codegen-service/templates/modules/shop/module.json`
- Create: `codegen-service/templates/modules/activity/module.json`
- Create: `codegen-service/scripts/injector.js`
- Create: `codegen-service/tests/injector.test.js`

- [ ] **Step 1: 创建 module.json 规范示例 (cms)**

```json
{
  "name": "cms",
  "displayName": "内容管理",
  "version": "1.0.0",
  "description": "文章发布、公告管理、轮播图",
  "dependencies": ["base"],
  "pages": [
    { "path": "pages/article/list", "style": {} },
    { "path": "pages/article/detail", "style": {} }
  ],
  "tabbar": null,
  "api": ["content.js"],
  "components": [],
  "configSchema": {
    "categories": { "type": "array", "default": ["健康资讯", "育儿知识"] }
  }
}
```

- [ ] **Step 2: 创建 injector.js 测试**

```javascript
const { injectConfig, replaceBrand } = require('../scripts/injector');

describe('Config Injector', () => {
  test('injects brand config into template content', () => {
    const template = 'const APP_NAME = "{{appName}}";\nconst COLOR = "{{primaryColor}}";';
    const brand = { appName: '宝宝健康助手', primaryColor: '#ff5722' };
    const result = injectConfig(template, brand);
    expect(result).toContain('宝宝健康助手');
    expect(result).toContain('#ff5722');
    expect(result).not.toContain('{{appName}}');
  });

  test('injects API endpoint and project credentials', () => {
    const template = "const BASE_URL = '{{baasBaseURL}}';\nconst API_KEY = '{{apiKey}}';";
    const config = { baasBaseURL: 'https://api.example.com', apiKey: 'ak_test123' };
    const result = injectConfig(template, config);
    expect(result).toContain('https://api.example.com');
    expect(result).toContain('ak_test123');
  });

  test('replaces all occurrences of a variable', () => {
    const template = '{{color}} is the color. Use {{color}} everywhere.';
    const result = injectConfig(template, { color: 'red' });
    expect(result).toBe('red is the color. Use red everywhere.');
  });

  test('leaves non-matching template tags intact', () => {
    const template = '{{appName}} {{nonexistent}}';
    const result = injectConfig(template, { appName: 'Test' });
    expect(result).toBe('Test {{nonexistent}}');
  });
});
```

- [ ] **Step 3: 实现 scripts/injector.js**

```javascript
/**
 * 将配置值注入模板内容，替换所有 {{key}} 占位符
 */
function injectConfig(template, config) {
  let result = template;
  for (const [key, value] of Object.entries(config)) {
    const regex = new RegExp(`\\{\\{${key}\\}\\}`, 'g');
    result = result.replace(regex, String(value));
  }
  return result;
}

/**
 * 生成完整的注入配置对象
 */
function buildInjectionConfig(project) {
  return {
    appName: project.brand_config?.appName || project.name,
    primaryColor: project.brand_config?.primaryColor || '#4caf50',
    secondaryColor: project.brand_config?.secondaryColor || '#ff9800',
    logo: project.brand_config?.logo || '/static/logo.png',
    footer: project.brand_config?.footer || `© ${new Date().getFullYear()} ${project.name}`,
    baasBaseURL: process.env.BAAS_BASE_URL || 'https://api.hxbaby.com',
    apiKey: project.api_key,
    projectId: String(project.id),
    wxAppId: project.wx_app_id || 'YOUR_WX_APP_ID',
  };
}

module.exports = { injectConfig, buildInjectionConfig };
```

- [ ] **Step 4: 运行测试**

```bash
npm test
# 预期: 全部 PASS
```

- [ ] **Step 5: Commit**

---

### Task 6: 模板编排 + ZIP打包器

**Files:**
- Create: `codegen-service/scripts/compose.js`
- Create: `codegen-service/scripts/packager.js`
- Create: `codegen-service/tests/compose.test.js`

- [ ] **Step 1: 实现 compose.js (核心编排)**

```javascript
const fs = require('fs');
const path = require('path');
const Handlebars = require('handlebars');
const { resolveModules, detectConflicts } = require('./resolver');
const { injectConfig, buildInjectionConfig } = require('./injector');

const TEMPLATES_DIR = path.join(__dirname, '..', 'templates');

/**
 * 编排主函数：根据项目配置生成完整的 UniApp 源码
 */
async function compose(project) {
  const warnings = [];
  const modules = resolveModules(project.modules);

  // 检查冲突
  const conflictWarnings = detectConflicts(modules);
  warnings.push(...conflictWarnings);

  // 准备输出目录
  const outputDir = path.join(__dirname, '..', 'output', project.build_task_id);
  fs.mkdirSync(outputDir, { recursive: true });

  const injectionConfig = buildInjectionConfig(project);

  // Step 1: 复制基础骨架
  const baseDir = path.join(TEMPLATES_DIR, 'base');
  copyAndInjectDir(baseDir, outputDir, injectionConfig);

  // Step 2: 合并 pages.json 路由
  const pagesJSON = mergePagesJSON(modules, outputDir);
  fs.writeFileSync(
    path.join(outputDir, 'pages.json'),
    JSON.stringify(pagesJSON, null, 2)
  );

  // Step 3: 复制各模块页面
  for (const mod of modules) {
    if (mod === 'base') continue;
    const modDir = path.join(TEMPLATES_DIR, 'modules', mod);
    if (!fs.existsSync(modDir)) {
      warnings.push(`模块 "${mod}" 模板不存在，已跳过`);
      continue;
    }
    const modOutputDir = path.join(outputDir, 'src');
    copyAndInjectDir(modDir, modOutputDir, injectionConfig);
  }

  // Step 4: 生成 package.json
  const pkgJSON = generatePackageJSON(project);
  fs.writeFileSync(
    path.join(outputDir, 'package.json'),
    JSON.stringify(pkgJSON, null, 2)
  );

  // Step 5: 生成 README
  const readme = generateREADME(project, modules, warnings);
  fs.writeFileSync(path.join(outputDir, 'README.md'), readme);

  return { outputDir, warnings };
}

/**
 * 合并所有模块的 pages.json 路由
 */
function mergePagesJSON(modules, outputDir) {
  let pages = [];
  let subPackages = {};

  for (const mod of modules) {
    const moduleJSONPath = path.join(TEMPLATES_DIR, 'modules', mod, 'module.json');
    if (!fs.existsSync(moduleJSONPath)) continue;

    const modConfig = JSON.parse(fs.readFileSync(moduleJSONPath, 'utf-8'));
    if (modConfig.pages) {
      for (const p of modConfig.pages) {
        pages.push({ path: p.path, style: p.style || {} });
      }
    }
  }

  return {
    pages,
    subPackages,
    globalStyle: {
      navigationBarTextStyle: 'black',
      navigationBarTitleText: '{{appName}}',
      navigationBarBackgroundColor: '#ffffff',
    },
    tabBar: buildTabBar(modules),
  };
}

/**
 * 根据选中模块构建 TabBar
 */
function buildTabBar(modules) {
  const tabItems = [];
  const tabOrder = ['cms', 'ai-advisor', 'shop', 'activity', 'member'];

  for (const mod of tabOrder) {
    if (!modules.includes(mod)) continue;
    const modJSON = path.join(TEMPLATES_DIR, 'modules', mod, 'module.json');
    if (!fs.existsSync(modJSON)) continue;

    const config = JSON.parse(fs.readFileSync(modJSON, 'utf-8'));
    if (config.tabbar) {
      tabItems.push(config.tabbar);
    }
  }

  if (tabItems.length === 0) return null;
  return { list: tabItems.slice(0, 5), color: '#999', selectedColor: '{{primaryColor}}' };
}

function copyAndInjectDir(srcDir, destDir, config) {
  if (!fs.existsSync(srcDir)) return;
  const entries = fs.readdirSync(srcDir, { withFileTypes: true });

  for (const entry of entries) {
    const srcPath = path.join(srcDir, entry.name);
    const destPath = path.join(destDir, entry.name.replace('.hbs', ''));

    if (entry.isDirectory()) {
      fs.mkdirSync(destPath, { recursive: true });
      copyAndInjectDir(srcPath, destPath, config);
    } else {
      let content = fs.readFileSync(srcPath, 'utf-8');
      // 只有 .hbs 文件做变量注入
      if (entry.name.endsWith('.hbs')) {
        content = injectConfig(content, config);
      }
      fs.writeFileSync(destPath, content);
    }
  }
}

function generatePackageJSON(project) {
  return {
    name: project.name.replace(/\s+/g, '-').toLowerCase(),
    version: '1.0.0',
    scripts: { 'dev:mp-weixin': 'uni -p mp-weixin', 'build:mp-weixin': 'uni build -p mp-weixin' },
    dependencies: {
      'vue': '^3.4.0',
      'uni-app': '^3.0.0',
      'pinia': '^2.1.0',
    },
  };
}

function generateREADME(project, modules, warnings) {
  return `# ${project.name}

## 模块列表
${modules.map(m => `- [x] ${m}`).join('\n')}

## 快速开始

1. 安装依赖: \`npm install\`
2. 配置微信AppID: 在 \`src/manifest.json\` 中填写微信小程序AppID
3. 启动开发: \`npm run dev:mp-weixin\`
4. 打开微信开发者工具，导入 \`dist/dev/mp-weixin\` 目录

## API配置
- BaaS地址: ${process.env.BAAS_BASE_URL || 'https://api.hxbaby.com'}
- API Key: ${project.api_key}

## 接入指引
详细接入指引见管理后台 → 项目 → 接入指引

${warnings.length > 0 ? `## 注意事项\n${warnings.map(w => `- ⚠️ ${w}`).join('\n')}` : ''}
`;
}

module.exports = { compose, mergePagesJSON };
```

- [ ] **Step 2: 实现 packager.js**

```javascript
const fs = require('fs');
const path = require('path');
const archiver = require('archiver');
const crypto = require('crypto');

/**
 * 将输出目录打包为 ZIP 并返回文件信息
 */
async function packageToZip(outputDir) {
  return new Promise((resolve, reject) => {
    const zipPath = outputDir + '.zip';
    const output = fs.createWriteStream(zipPath);
    const archive = archiver('zip', { zlib: { level: 9 } });

    output.on('close', () => {
      const fileBuffer = fs.readFileSync(zipPath);
      const md5 = crypto.createHash('md5').update(fileBuffer).digest('hex');
      const size = fs.statSync(zipPath).size;
      resolve({ zipPath, md5, size });
    });

    archive.on('error', reject);
    archive.pipe(output);
    archive.directory(outputDir, path.basename(outputDir));
    archive.finalize();
  });
}

module.exports = { packageToZip };
```

- [ ] **Step 3: 创建 server.js (Express API)**

```javascript
const express = require('express');
const { compose } = require('./scripts/compose');
const { packageToZip } = require('./scripts/packager');
const { v4: uuidv4 } = require('uuid');

const app = express();
app.use(express.json());

// 健康检查
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'codegen-service' });
});

// 获取可用模块列表
app.get('/api/modules', (req, res) => {
  const fs = require('fs');
  const path = require('path');
  const modulesDir = path.join(__dirname, 'templates', 'modules');
  const modules = fs.readdirSync(modulesDir)
    .filter(d => fs.statSync(path.join(modulesDir, d)).isDirectory())
    .map(d => {
      const configPath = path.join(modulesDir, d, 'module.json');
      if (fs.existsSync(configPath)) {
        return JSON.parse(fs.readFileSync(configPath, 'utf-8'));
      }
      return { name: d, displayName: d };
    });
  res.json({ modules });
});

// 代码生成 API
app.post('/api/build', async (req, res) => {
  try {
    const { project } = req.body;
    if (!project || !project.modules) {
      return res.status(400).json({ error: '缺少必要参数' });
    }

    project.build_task_id = uuidv4();

    const { outputDir, warnings } = await compose(project);
    const { zipPath, md5, size } = await packageToZip(outputDir);

    res.json({
      task_id: project.build_task_id,
      status: 'done',
      output_dir: outputDir,
      zip_path: zipPath,
      md5,
      size_bytes: size,
      warnings,
    });
  } catch (err) {
    res.status(500).json({ status: 'failed', error: err.message });
  }
});

const PORT = process.env.PORT || 3002;
app.listen(PORT, () => {
  console.log(`CodeGen service running on port ${PORT}`);
});
```

- [ ] **Step 4: 运行测试**

```bash
npm test
# 预期: 全部 PASS
```

- [ ] **Step 5: Commit**

```bash
git add codegen-service/scripts/compose.js codegen-service/scripts/packager.js codegen-service/server.js
git commit -m "feat: add code generation compose engine and Express API"
```

---

## Phase 3: BaaS运行时API (Week 3-4)

### Task 7: BaaS认证中间件 + 内容API

**Files:**
- Create: `biz-service/internal/middleware/baas_auth.go`
- Create: `biz-service/internal/middleware/project_isolate.go`
- Create: `biz-service/internal/model/cms_article.go`
- Create: `biz-service/internal/repository/cms_article.go`
- Create: `biz-service/internal/handler/baas_content.go`

- [ ] **Step 1: 实现 BaaS API Key 认证中间件**

```go
package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/hxbaby/biz-service/internal/model"
)

// BaaSAuth 验证来自小程序的API Key
func BaaSAuth(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" {
            apiKey = c.Query("api_key")
        }
        if apiKey == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少API Key"})
            c.Abort()
            return
        }

        var project model.MiniappProject
        if err := db.Where("api_key = ? AND status = ?", apiKey, "active").First(&project).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的API Key"})
            c.Abort()
            return
        }

        c.Set("project_id", project.ID)
        c.Set("customer_id", project.CustomerID)
        c.Next()
    }
}

// ProjectIsolate 确保数据隔离在 project 级别
func ProjectIsolate() gin.HandlerFunc {
    return func(c *gin.Context) {
        projectID, exists := c.Get("project_id")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "项目信息缺失"})
            c.Abort()
            return
        }
        c.Set("project_id", projectID)
        c.Next()
    }
}
```

- [ ] **Step 2: 实现 CmsArticle Model + Repository**

```go
// model/cms_article.go
package model

import "time"

type CmsArticle struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    ProjectID  uint      `gorm:"index;not null" json:"project_id"`
    Title      string    `gorm:"size:200;not null" json:"title"`
    Content    string    `gorm:"type:text" json:"content"`
    Summary    string    `gorm:"size:500" json:"summary,omitempty"`
    CoverImage string    `gorm:"size:500" json:"cover_image,omitempty"`
    Category   string    `gorm:"size:50" json:"category,omitempty"`
    Tags       StringSlice `gorm:"type:jsonb;default:'[]'" json:"tags"`
    AIGenerated bool     `gorm:"default:false" json:"ai_generated"`
    AIPrompt   string    `gorm:"type:text" json:"ai_prompt,omitempty"`
    IsPublished bool     `gorm:"default:false" json:"is_published"`
    ViewCount   int      `gorm:"default:0" json:"view_count"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

- [ ] **Step 3: 实现 BaaS Content Handler**

```go
// handler/baas_content.go
package handler

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/hxbaby/biz-service/internal/model"
    "github.com/hxbaby/biz-service/pkg/response"
)

type BaaSContentHandler struct {
    db *gorm.DB
}

func NewBaaSContentHandler(db *gorm.DB) *BaaSContentHandler {
    return &BaaSContentHandler{db: db}
}

func (h *BaaSContentHandler) ListArticles(c *gin.Context) {
    projectID := c.GetUint("project_id")
    category := c.Query("category")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

    var articles []model.CmsArticle
    query := h.db.Where("project_id = ? AND is_published = ?", projectID, true)
    if category != "" {
        query = query.Where("category = ?", category)
    }

    var total int64
    query.Model(&model.CmsArticle{}).Count(&total)
    query.Order("created_at DESC").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&articles)

    response.OK(c, gin.H{
        "items":     articles,
        "total":     total,
        "page":      page,
        "page_size": pageSize,
    })
}

func (h *BaaSContentHandler) GetArticle(c *gin.Context) {
    projectID := c.GetUint("project_id")
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

    var article model.CmsArticle
    if err := h.db.Where("id = ? AND project_id = ?", id, projectID).First(&article).Error; err != nil {
        response.Error(c, http.StatusNotFound, "文章不存在")
        return
    }

    // 增加阅读量
    h.db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
    article.ViewCount++

    response.OK(c, article)
}
```

- [ ] **Step 4: 注册 BaaS 路由**

```go
// router.go (baas 路由组)
baas := r.Group("/baas/v1")
baas.Use(middleware.BaaSAuth(config.DB), middleware.ProjectIsolate())

contentH := handler.NewBaaSContentHandler(config.DB)
baas.GET("/articles", contentH.ListArticles)
baas.GET("/articles/:id", contentH.GetArticle)
```

- [ ] **Step 5: Commit**

---

### Task 8: AI Bridge + AI增强API (文案/海报)

**Files:**
- Create: `biz-service/internal/service/ai_bridge.go`
- Create: `biz-service/internal/handler/ai_enhance.go`

- [ ] **Step 1: 实现 AI Bridge service/ai_bridge.go**

```go
package service

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type AIBridge struct {
    baseURL    string
    httpClient *http.Client
}

type AIGenerateRequest struct {
    Prompt     string                 `json:"prompt"`
    Model      string                 `json:"model,omitempty"`
    MaxTokens  int                    `json:"max_tokens,omitempty"`
    Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type AIGenerateResponse struct {
    Content string `json:"content"`
    Tokens  int    `json:"tokens_used"`
    Model   string `json:"model"`
}

func NewAIBridge(baseURL string) *AIBridge {
    return &AIBridge{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 120 * time.Second,
        },
    }
}

func (b *AIBridge) GenerateText(prompt string, maxTokens int) (*AIGenerateResponse, error) {
    req := AIGenerateRequest{
        Prompt:    prompt,
        MaxTokens: maxTokens,
    }

    body, _ := json.Marshal(req)
    resp, err := b.httpClient.Post(
        b.baseURL+"/ai/generate",
        "application/json",
        bytes.NewReader(body),
    )
    if err != nil {
        return nil, fmt.Errorf("AI服务不可用: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("AI服务返回错误: %s", string(bodyBytes))
    }

    var result AIGenerateResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}

// GenerateArticle 生成文章
func (b *AIBridge) GenerateArticle(topic, category string) (*AIGenerateResponse, error) {
    prompt := fmt.Sprintf(`你是一个专业的儿童健康科普作者。
请根据以下主题写一篇公众号文章（约800字）：
主题：%s
分类：%s
要求：专业易懂、适合宝妈阅读、结构清晰有小标题。`, topic, category)
    return b.GenerateText(prompt, 2000)
}

// GenerateSummary 生成摘要
func (b *AIBridge) GenerateSummary(article string) (*AIGenerateResponse, error) {
    prompt := fmt.Sprintf("请为以下文章生成一段约100字的摘要：\n\n%s", article)
    return b.GenerateText(prompt, 200)
}

// GenerateActivityCopy 生成活动文案
func (b *AIBridge) GenerateActivityCopy(title, description string) (*AIGenerateResponse, error) {
    prompt := fmt.Sprintf(`你是一个活动策划专家。请为以下活动生成营销文案：
活动标题：%s
活动描述：%s
请生成：
1. 朋友圈推广文案（150字以内）
2. 活动详情页文案（300字）
3. 群发通知文案（100字）`, title, description)
    return b.GenerateText(prompt, 1500)
}

// GenerateSellingPoints 提炼卖点
func (b *AIBridge) GenerateSellingPoints(productName, productDesc string) (*AIGenerateResponse, error) {
    prompt := fmt.Sprintf(`你是一个母婴产品营销专家。请为以下产品提炼核心卖点：
产品名称：%s
产品描述：%s
要求：3-5个卖点，每个卖点20字以内，突出对宝宝/宝妈的价值。`, productName, productDesc)
    return b.GenerateText(prompt, 500)
}

// GenerateActivityReport 生成活动复盘报告
func (b *AIBridge) GenerateActivityReport(activityName string, stats map[string]interface{}) (*AIGenerateResponse, error) {
    prompt := fmt.Sprintf(`你是一个数据分析专家。请根据以下活动数据生成复盘报告：
活动名称：%s
活动数据：%v
请分析：参与情况、转化效果、亮点与不足、改进建议。`, activityName, stats)
    return b.GenerateText(prompt, 1500)
}
```

- [ ] **Step 2: 实现 AI增强 Handler**

```go
// handler/ai_enhance.go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/hxbaby/biz-service/internal/service"
    "github.com/hxbaby/biz-service/pkg/response"
)

type AIEnhanceHandler struct {
    bridge *service.AIBridge
}

func NewAIEnhanceHandler(bridge *service.AIBridge) *AIEnhanceHandler {
    return &AIEnhanceHandler{bridge: bridge}
}

func (h *AIEnhanceHandler) GenerateArticle(c *gin.Context) {
    var req struct {
        Topic    string `json:"topic" binding:"required"`
        Category string `json:"category"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    result, err := h.bridge.GenerateArticle(req.Topic, req.Category)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateSummary(c *gin.Context) {
    var req struct {
        Content string `json:"content" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    result, err := h.bridge.GenerateSummary(req.Content)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateActivityCopy(c *gin.Context) {
    var req struct {
        Title       string `json:"title" binding:"required"`
        Description string `json:"description"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    result, err := h.bridge.GenerateActivityCopy(req.Title, req.Description)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateSellingPoints(c *gin.Context) {
    var req struct {
        Name        string `json:"name" binding:"required"`
        Description string `json:"description"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    result, err := h.bridge.GenerateSellingPoints(req.Name, req.Description)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, result)
}

func (h *AIEnhanceHandler) GenerateActivityReport(c *gin.Context) {
    var req struct {
        ActivityName string                 `json:"activity_name" binding:"required"`
        Stats        map[string]interface{} `json:"stats"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    result, err := h.bridge.GenerateActivityReport(req.ActivityName, req.Stats)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.OK(c, result)
}
```

- [ ] **Step 3: 注册AI增强路由**

```go
aiBridge := service.NewAIBridge(cfg.AIServiceURL)
aiEnhanceH := handler.NewAIEnhanceHandler(aiBridge)

// AI增强接口（需管理端登录态）
enhance := v1.Group("/ai", middleware.AuthRequired(cfg.JWTSecret))
enhance.POST("/generate-article", aiEnhanceH.GenerateArticle)
enhance.POST("/generate-summary", aiEnhanceH.GenerateSummary)
enhance.POST("/generate-activity-copy", aiEnhanceH.GenerateActivityCopy)
enhance.POST("/generate-selling-points", aiEnhanceH.GenerateSellingPoints)
enhance.POST("/activity-report", aiEnhanceH.GenerateActivityReport)
```

- [ ] **Step 4: Commit**

---

## Phase 4-5: React管理前端 + 集成联调 (Week 5-8)

### Task 9: React管理后台骨架 + 登录注册

**Files:**
- Create: `factory-admin/` (Vite + React + Ant Design 5 项目)
- Create: 登录/注册页面、Dashboard、路由配置

- [ ] **Step 1: 初始化项目**

```bash
npm create vite@latest factory-admin -- --template react-ts
cd factory-admin
npm install antd @ant-design/icons axios react-router-dom
npm install --save-dev @types/react-router-dom
```

- [ ] **Step 2: 实现 axios client (src/api/client.ts)**

```typescript
import axios from 'axios';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 30000,
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(err);
  }
);

export default client;
```

- [ ] **Step 3: 实现核心页面**

详细代码略（登录、注册、Dashboard、项目列表、创建向导、模块选择器、品牌配置、构建进度、接入指引、计费页面），遵循 React + Ant Design 5 最佳实践。

- [ ] **Step 4: Commit**

---

### Task 10: 构建触发 + BaaS 电商/活动/用户 API

**Files:**
（Go层补充电商、活动、用户、预约等BaaS Handler + Repository；构建状态回调；前端构建触发 + 进度轮询）

按同一模式实现，代码结构同 Task 7。

---

### Task 11: 端到端集成测试 + 一键部署

**Files:**
- Create: `docker-compose.full.yml` (包含所有服务)
- Create: `scripts/deploy.sh`
- Create: 集成测试脚本

验证全链路：客户注册 → 创建项目 → 勾选模块 → 一键构建 → 下载源码 → 微信开发者工具编译 → BaaS API 调用 → AI对话返回。

---

## 自检清单

- [x] 每个核心Task都有测试→实现→验证→Commit流程
- [x] resolver/injector/compose 三组件接口一致（project.modules → resolver → injector → compose）
- [x] Go model 字段与设计文档数据模型匹配
- [x] BaaS API 鉴权通过 X-API-Key header + project_id 隔离
- [x] AI Bridge 统一调用 Python AI 服务，不在各业务Handler中散落调用
- [x] 模块拓扑排序确保依赖模块先于被依赖模块加载
