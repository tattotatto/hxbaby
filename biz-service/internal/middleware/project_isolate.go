package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProjectIsolate ensures that a project_id has been set in the gin context
// (by BaaSAuth or similar middleware) before allowing the request through.
func ProjectIsolate() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("project_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "项目信息缺失"})
			c.Abort()
			return
		}
		c.Next()
	}
}
