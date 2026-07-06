package middleware

import (
	"net"
	"strings"

	"forum-gateway/handler"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// WebhookGuard 限制 webhook 回调只能被审核服务和本地访问（本地方便测试）
// 生产环境：审核服务在 webhook URL 中带上 ?token=xxx 参数
// 本地环境：localhost 直连放行
func WebhookGuard() gin.HandlerFunc {
	expectedToken := viper.GetString("audit.webhook_token")

	return func(c *gin.Context) {
		if isLocalhost(c.Request.Host) || isLocalhost(c.ClientIP()) {
			c.Next()
			return
		}

		// 生产环境：校验 token（优先从 query 取，其次从 X-Webhook-Token 头取）
		token := c.Query("token")
		if token == expectedToken {
			c.Next()
			return
		}

		handler.SendError(c, errno.ErrPermissionDenied, nil, "webhook 回调认证失败", handler.GetLine())
		c.Abort()
	}
}

func isLocalhost(addr string) bool {
	if h, _, err := net.SplitHostPort(addr); err == nil {
		addr = h
	}
	addr = strings.ToLower(addr)
	switch addr {
	case "localhost", "127.0.0.1", "::1":
		return true
	}
	if ip := net.ParseIP(addr); ip != nil && ip.IsLoopback() {
		return true
	}
	return false
}
