package middleware

import (
	"forum-gateway/handler"
	"forum-gateway/pkg/auth"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware ... 认证中间件
// limit 为限制的权限等级
func AuthMiddleware(limit uint32) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		ctx, err := auth.ParseRequest(c)
		if err != nil {
			handler.SendResponse(c, errno.ErrAuthToken, err.Error())
			c.Abort()
			return
		} else if ctx.Role&limit == 0 {
			handler.SendResponse(c, errno.ErrPermissionDenied, "")
			c.Abort()
			return
		}

		c.Set("userId", ctx.Id)
		c.Set("role", ctx.Role)
		c.Set("expiresAt", ctx.ExpiresAt)

		c.Next()
	}
}
