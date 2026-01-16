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

	v1.Use(middleware.RequestId())

	// 公开接口
	v1.POST("/register", userHandler.Register)
	v1.POST("/login", userHandler.Login)

	// 需要认证的接口（创建一个新的作用域，Use() 方法会将中间件应用到后续注册的所有路由上）
	v1.Use(middleware.AuthorizationCheck(), middleware.PermissionCheck())

	v1.POST("/users/list", userHandler.GetUserList)
	v1.POST("/users/get", userHandler.GetUserByID)
	v1.POST("/users/update", userHandler.UpdateUser)
	v1.POST("/users/delete", userHandler.DeleteUser)

	return router
}
