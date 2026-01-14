package middleware

import (
	"net/http"
	"strings"
	"users-by-go-example/global"
	"users-by-go-example/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查是否在白名单中
		cfg := global.GetConfig()
		path := ctx.Request.URL.Path
		for _, whitePath := range cfg.WhiteList {
			if path == whitePath {
				ctx.Next()
				return
			}
		}

		// 获取 Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "请提供认证令牌",
			})
			ctx.Abort()
			return
		}

		// 检查 Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "认证令牌格式错误",
			})
			ctx.Abort()
			return
		}

		// 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "认证令牌无效或已过期",
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("userId", claims.UserID)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
