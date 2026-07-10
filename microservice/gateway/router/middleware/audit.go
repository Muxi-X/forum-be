package middleware

import (
	"strings"

	"forum-gateway/handler/audit"

	"github.com/gin-gonic/gin"
)

// Audit 按请求路径+方法分发到审核流程
// 客户端设置 X-Request-Audit: true 头即可触发先审后发
func Audit(auditApi *audit.Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Request-Audit") != "true" {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		switch {
		// post
		case strings.Contains(path, "/api/v1/post") && method == "POST":
			auditApi.Post(c)
		case strings.Contains(path, "/api/v1/post") && method == "PUT":
			auditApi.PostUpdate(c)

		// sip-score
		case strings.Contains(path, "/api/v1/sip-score") && method == "POST":
			auditApi.CreateSipScore(c)
		case strings.Contains(path, "/api/v1/sip-score") && method == "PUT":
			auditApi.UpdateSipScore(c)

		default:
			c.Next()
			return
		}
		c.Abort()
	}
}
