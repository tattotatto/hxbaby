package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hxbaby/biz-service/pkg/response"
)

func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("tenant_id")
		if !exists {
			response.Error(c, http.StatusForbidden, "租户信息缺失")
			c.Abort()
			return
		}
		c.Next()
	}
}
