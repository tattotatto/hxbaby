package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/hxbaby/biz-service/internal/model"
)

var DB *gorm.DB

type Config struct {
	ServerPort     string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	AIServiceURL   string
	CodegenURL     string
	Environment    string
}

func Load() *Config {
	return &Config{
		ServerPort:   getEnv("BIZ_SERVICE_PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgresql://hxbaby:hxbaby_dev@localhost:5432/hxbaby?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:    getEnv("JWT_SECRET", "dev-secret-change-me-in-production"),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8001"),
		CodegenURL:   getEnv("CODEGEN_SERVICE_URL", "http://localhost:3002"),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}
}

func InitDB(databaseURL string) error {
	var err error
	logLevel := logger.Info
	if getEnv("ENVIRONMENT", "development") == "production" {
		logLevel = logger.Warn
	}

	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	log.Println("Database connected")
	return nil
}

func AutoMigrate() {
	if DB == nil {
		log.Println("DB not initialized, skipping migration")
		return
	}
	DB.AutoMigrate(
		&model.Tenant{},
		&model.User{},
		&model.Child{},
		&model.Conversation{},
		&model.Message{},
		&model.ProductCategory{},
		&model.Product{},
		&model.Assessment{},
		&model.Recommendation{},
		&model.Customer{},
		&model.MiniappProject{},
		&model.BuildTask{},
		&model.CmsArticle{},
		&model.BProduct{},
		&model.BOrder{},
		&model.BActivity{},
		&model.BActivitySignup{},
		&model.BUser{},
		&model.BBooking{},
	)
	log.Println("Database migration completed")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
