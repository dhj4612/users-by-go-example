package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username" json:"username"`
	Password   string    `gorm:"column:password;type:varchar(255);not null" json:"-"` // json:"-" 表示不返回密码
	NikeName   string    `gorm:"column:nike_name;type:varchar(50)" json:"nikeName"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Delete     int       `gorm:"column:delete;type:tinyint(1);default:0" json:"-"` // 0-未删除 1-已删除
}

// TableName 指定表名
func (*User) TableName() string {
	return "users"
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	NikeName string `json:"nikeName" binding:"max=50"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       int64  `json:"id" binding:"required"`
	NikeName string `json:"nikeName" binding:"max=50"`
	Password string `json:"password" binding:"omitempty,min=6,max=50"`
}

// GetUserListRequest 获取用户列表请求
type GetUserListRequest struct {
	Page     int `json:"page" binding:"omitempty,min=1"`
	PageSize int `json:"pageSize" binding:"omitempty,min=1,max=100"`
}

// GetUserByIDRequest 根据 ID 获取用户请求
type GetUserByIDRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	NikeName   string    `json:"nikeName"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// maskUsername 用户名脱敏
// 规则：保留首尾各1个字符，中间用 * 替代
// 例如：testuser -> t******r, abc -> a*c, ab -> a*
func maskUsername(username string) string {
	length := len(username)
	if length == 0 {
		return ""
	}
	if length == 1 {
		return username
	}
	if length == 2 {
		return string(username[0]) + "*"
	}
	// 长度大于等于3，保留首尾各1个字符
	masked := string(username[0])
	for i := 1; i < length-1; i++ {
		masked += "*"
	}
	masked += string(username[length-1])
	return masked
}

// ToResponse 转换为响应对象（用户名脱敏）
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   maskUsername(u.Username),
		NikeName:   u.NikeName,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}
