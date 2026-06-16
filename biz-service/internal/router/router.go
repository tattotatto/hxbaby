package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/config"
	"github.com/hxbaby/biz-service/internal/handler"
	"github.com/hxbaby/biz-service/internal/middleware"
	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/internal/service"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "biz-service",
			"version": "0.1.0",
		})
	})

	// 初始化依赖
	userRepo := repository.NewUserRepo(config.DB)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	authH := handler.NewAuthHandler(authSvc)

	// 初始化
	childH := handler.NewChildHandler(repository.NewChildRepo(config.DB))
	convRepo := repository.NewConversationRepo(config.DB)
	msgRepo := repository.NewMessageRepo(config.DB)
	productRepo := repository.NewProductRepo(config.DB)

	aiClient := service.NewAIClient(cfg.AIServiceURL)
	chatH := handler.NewChatHandler(aiClient)

	// 公开路由
	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", authH.Login)
	v1.POST("/auth/register", authH.Register)

	// 认证+租户路由
	auth := v1.Group("", middleware.AuthRequired(cfg.JWTSecret), middleware.TenantIsolation())
	auth.POST("/conversations/:id/chat", chatH.Chat)
	auth.GET("/children", childH.List)
	auth.POST("/children", childH.Create)
	auth.GET("/children/:id", childH.Get)
	auth.POST("/children/:id/growth", childH.AddGrowth)
	auth.GET("/products", handler.NewProductHandler(productRepo).List)
	auth.GET("/products/match", handler.NewProductHandler(productRepo).Match)

	// Factory customer routes
	customerRepo := repository.NewCustomerRepo(config.DB)
	customerSvc := service.NewCustomerService(customerRepo, cfg.JWTSecret)
	customerH := handler.NewCustomerHandler(customerSvc)

	factoryAuth := v1.Group("/factory/auth")
	factoryAuth.POST("/register", customerH.Register)
	factoryAuth.POST("/login", customerH.Login)

	// Factory project routes (require customer JWT)
	projectRepo := repository.NewProjectRepo(config.DB)
	projectSvc := service.NewProjectService(projectRepo, customerRepo)
	projectH := handler.NewProjectHandler(projectSvc)

	factoryAuthRequired := v1.Group("/factory", middleware.AuthRequired(cfg.JWTSecret))
	factoryAuthRequired.GET("/projects", projectH.List)
	factoryAuthRequired.POST("/projects", projectH.Create)
	factoryAuthRequired.GET("/projects/:id", projectH.Get)
	factoryAuthRequired.PUT("/projects/:id", projectH.Update)

	// Build routes (JWT protected)
	buildRepo := repository.NewBuildTaskRepo(config.DB)
	buildH := handler.NewBuildHandler(buildRepo, projectRepo, cfg.CodegenURL)

	factoryAuthRequired.POST("/projects/:id/build", buildH.TriggerBuild)
	factoryAuthRequired.GET("/projects/:id/builds", buildH.GetBuildHistory)
	v1.GET("/builds/:id/status", buildH.GetBuildStatus)
	v1.GET("/builds/:id/download", buildH.DownloadBuild)

	// BaaS routes (X-API-Key auth for mini-programs)
	baas := r.Group("/baas/v1")
	baas.Use(middleware.BaaSAuth(config.DB), middleware.ProjectIsolate())

	contentH := handler.NewBaaSContentHandler(config.DB)
	baas.GET("/articles", contentH.ListArticles)
	baas.GET("/articles/:id", contentH.GetArticle)

	// BaaS shop routes
	shopH := handler.NewBaaSShopHandler(config.DB)
	baas.GET("/products", shopH.ListProducts)
	baas.GET("/products/:id", shopH.GetProduct)
	baas.POST("/orders", shopH.CreateOrder)
	baas.GET("/orders", shopH.ListOrders)
	baas.GET("/orders/:id", shopH.GetOrder)

	// BaaS activity routes
	activityH := handler.NewBaaSActivityHandler(config.DB)
	baas.GET("/activities", activityH.ListActivities)
	baas.GET("/activities/:id", activityH.GetActivity)
	baas.POST("/activities/:id/signup", activityH.SignupActivity)
	baas.POST("/activities/:id/checkin", activityH.CheckinActivity)

	// BaaS user routes
	userH := handler.NewBaaSUserHandler(config.DB)
	baas.POST("/auth/wx-login", userH.WxLogin)
	baas.GET("/user/profile", userH.GetProfile)
	baas.PUT("/user/profile", userH.UpdateProfile)

	// BaaS booking routes
	bookingH := handler.NewBaaSBookingHandler(config.DB)
	baas.GET("/bookings/slots", bookingH.ListSlots)
	baas.POST("/bookings", bookingH.Create)
	baas.GET("/bookings", bookingH.List)

	// AI Enhance routes (require customer JWT)
	aiBridge := service.NewAIBridge(cfg.AIServiceURL)
	aiEnhanceH := handler.NewAIEnhanceHandler(aiBridge)

	aiEnhance := v1.Group("/ai", middleware.AuthRequired(cfg.JWTSecret))
	aiEnhance.POST("/generate-article", aiEnhanceH.GenerateArticle)
	aiEnhance.POST("/generate-summary", aiEnhanceH.GenerateSummary)
	aiEnhance.POST("/generate-activity-copy", aiEnhanceH.GenerateActivityCopy)
	aiEnhance.POST("/generate-selling-points", aiEnhanceH.GenerateSellingPoints)
	aiEnhance.POST("/activity-report", aiEnhanceH.GenerateActivityReport)

	// Knowledge base management (require customer JWT)
	knowledgeH := handler.NewKnowledgeHandler(cfg.AIServiceURL)
	knowledge := v1.Group("/knowledge", middleware.AuthRequired(cfg.JWTSecret))
	knowledge.POST("/upload", knowledgeH.Upload)
	knowledge.GET("/documents", knowledgeH.ListDocuments)
	knowledge.DELETE("/documents/:source", knowledgeH.DeleteDocument)
	knowledge.GET("/stats", knowledgeH.GetStats)

	_ = convRepo
	_ = msgRepo

	return r
}
