package handler

import (
	"net/http"
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

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PageResponse 分页响应
type PageResponse struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// 响应辅助函数
func success(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func badRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

func unauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

func notFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

func internalError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

// Register 用户注册
func (h *UserHandler) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	success(ctx, "注册成功", user)
}

// Login 用户登录
func (h *UserHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "参数错误: "+err.Error())
		return
	}

	token, err := h.userService.Login(&req)
	if err != nil {
		unauthorized(ctx, err.Error())
		return
	}

	success(ctx, "登录成功", gin.H{
		"token": token,
	})
}

// GetUserList 获取用户列表
func (h *UserHandler) GetUserList(ctx *gin.Context) {
	var req model.GetUserListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "参数错误: "+err.Error())
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
		internalError(ctx, "查询失败: "+err.Error())
		return
	}

	success(ctx, "查询成功", PageResponse{
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
		badRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.GetUserByID(req.ID)
	if err != nil {
		notFound(ctx, err.Error())
		return
	}

	success(ctx, "查询成功", user)
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var req model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.UpdateUser(req.ID, &req)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	success(ctx, "更新成功", user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	var req model.DeleteUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.DeleteUser(req.ID); err != nil {
		badRequest(ctx, err.Error())
		return
	}

	success(ctx, "删除成功", nil)
}
