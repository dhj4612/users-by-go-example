package handler

import (
	"users-by-go-example/internal/model"
	"users-by-go-example/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: &service.UserService{},
	}
}

// Register 用户注册
func (h *UserHandler) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	Success(ctx, "注册成功", user)
}

// Login 用户登录
func (h *UserHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	token, err := h.userService.Login(&req)
	if err != nil {
		Unauthorized(ctx, err.Error())
		return
	}

	Success(ctx, "登录成功", gin.H{
		"token": token,
	})
}

// GetUserList 获取用户列表
func (h *UserHandler) GetUserList(ctx *gin.Context) {
	var req model.GetUserListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	// 设置默认值
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	users, total, err := h.userService.GetUserList(page, pageSize)
	if err != nil {
		InternalError(ctx, "查询失败: "+err.Error())
		return
	}

	Success(ctx, "查询成功", PageResponse{
		List:     users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetUserByID 根据 ID 获取用户
func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	var req model.GetUserByIDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.GetUserByID(req.ID)
	if err != nil {
		NotFound(ctx, err.Error())
		return
	}

	Success(ctx, "查询成功", user)
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var req model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.UpdateUser(req.ID, &req)
	if err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	Success(ctx, "更新成功", user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	var req model.DeleteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.DeleteUser(req.ID); err != nil {
		BadRequest(ctx, err.Error())
		return
	}

	Success(ctx, "删除成功", nil)
}
