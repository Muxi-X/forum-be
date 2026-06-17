package middleware

import (
	"forum-gateway/handler"
	"forum-gateway/pkg/auth"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware ... 认证中间件（仅解析token，不校验权限等级）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, err := auth.ParseRequest(c)
		if err != nil {
			handler.SendError(c, errno.ErrAuthToken, nil, err.Error(), handler.GetLine())
			c.Abort()
			return
		}

		c.Set("userId", ctx.Id)
		c.Set("role", ctx.Role)
		c.Set("expiresAt", ctx.ExpiresAt)

		c.Next()
	}
}
