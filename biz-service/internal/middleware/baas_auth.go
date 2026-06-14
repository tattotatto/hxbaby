package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/hxbaby/biz-service/internal/model"
)

// BaaSAuth validates the X-API-Key header against the MiniappProject table
// and sets project_id and customer_id in the gin context for downstream handlers.
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
