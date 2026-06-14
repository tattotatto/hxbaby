package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/internal/config"
	"github.com/hxbaby/biz-service/internal/router"
)

func main() {
	cfg := config.Load()

	// 初始化数据库 (开发环境允许失败)
	if err := config.InitDB(cfg.DatabaseURL); err != nil {
		log.Printf("WARNING: Database not available: %v", err)
	} else {
		config.AutoMigrate()
	}

	// 设置 Gin 模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.Setup(cfg)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Biz service starting on %s (env: %s)", addr, cfg.Environment)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
