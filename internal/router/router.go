package router

import (
	"users-by-go-example/internal/handler"
	"users-by-go-example/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 创建用户处理器
	userHandler := handler.NewUserHandler()

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 公开接口（不需要认证）
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)

		// 需要认证的接口
		v1.Use(middleware.AuthMiddleware())
		{
			// 用户管理接口 - 统一使用 POST
			v1.POST("/users/list", userHandler.GetUserList)
			v1.POST("/users/get", userHandler.GetUserByID)
			v1.POST("/users/update", userHandler.UpdateUser)
			v1.POST("/users/delete", userHandler.DeleteUser)
		}
	}

	return router
}
