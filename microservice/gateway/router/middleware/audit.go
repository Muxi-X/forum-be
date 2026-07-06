package middleware

import (
	"strings"

	"forum-gateway/handler/audit"

	"github.com/gin-gonic/gin"
)

// Audit 按请求路径分发到审核流程
// 客户端设置 X-Request-Audit: true 头即可触发先审后发
func Audit(auditApi *audit.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Request-Audit") != "true" {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		switch {
		case strings.HasPrefix(path, "/api/v1/post"):
			auditApi.Post(c)
		// 以后扩展其他资源:
		// case strings.HasPrefix(path, "/api/v1/comment"):
		// 	auditApi.Comment(c)
		default:
			c.Next()
			return
		}
		c.Abort()
	}
}
