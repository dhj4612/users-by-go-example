package middleware

import (
	"net/http"
	"slices"
	"strings"
	"users-by-go-example/application"
	"users-by-go-example/logger"

	"github.com/gin-gonic/gin"
)

func PermissionCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		key := method + " " + path
		permitValue := application.GetApiPermitsMap()[key]

		var permits []string
		if strings.Contains(permitValue, ",") {
			permits = strings.Split(permitValue, ",")
		} else {
			permits = append(permits, permitValue)
		}

		log := logger.GetLogger(ctx)
		log.Info("Api=%s Permits=%s", key, permits)

		userId, exists := ctx.Get("userId")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "请提供认证令牌",
			})
			ctx.Abort()
			return
		}

		permitsOfUser := make([]string, 10)
		application.GetDB().Raw(`
						SELECT permit
						FROM permission p
						JOIN user_permission up ON p.id = up.permission_id
						WHERE up.user_id = ?`, userId).Scan(&permitsOfUser)

		log.Info("user permits=%v\n", permitsOfUser)

		if slices.Contains(permitsOfUser, "*") {
			ctx.Next()
			return
		}

		for i := 0; i < len(permits); i++ {
			allow := slices.Contains(permitsOfUser, permits[i])
			if !allow {
				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    http.StatusForbidden,
					"message": "未授权的访问",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
