package middleware

import (
	"strings"
	"users-by-go-example/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestId 中间件：从 header 中获取或生成 requestId，并创建 logger 存入 context
func RequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestId := ctx.GetHeader("x-request-id")
		if requestId == "" {
			requestId = strings.ReplaceAll(uuid.New().String(), "-", "")
		}

		log := logger.NewLogger(requestId)
		logger.SetLogger(ctx, log)

		ctx.Next()
	}
}
