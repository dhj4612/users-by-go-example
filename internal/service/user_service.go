package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users-by-go-example/global"
	"users-by-go-example/internal/model"
	"users-by-go-example/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// Register 用户注册
func (s *UserService) Register(req *model.RegisterRequest) (*model.UserResponse, error) {
	db := global.GetDB()
	rdb := global.GetRedis()
	ctx := context.Background()

	// 创建分布式锁，针对 username 加锁
	lock := utils.NewRedisLock(rdb, "register:"+req.Username, 10*time.Second)

	// 尝试获取锁，最多重试 3 次，每次间隔 100ms
	if err := lock.TryLock(ctx, 1, 100*time.Millisecond); err != nil {
		if errors.Is(err, utils.ErrLockFailed) {
			return nil, errors.New("系统繁忙，请稍后重试")
		}
		return nil, err
	}
	defer lock.Unlock(ctx)

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户名是否已存在
	var count int64
	if err := tx.Model(&model.User{}).Where("username = ? AND `delete` = 0", req.Username).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if count > 0 {
		tx.Rollback()
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		NikeName: req.NikeName,
		Delete:   0,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// Login 用户登录
func (s *UserService) Login(req *model.LoginRequest) (string, error) {
	db := global.GetDB()

	// 查询用户
	var user model.User
	if err := db.Where("username = ? AND `delete` = 0", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户名或密码错误")
		}
		return "", err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(page, pageSize int) ([]*model.UserResponse, int64, error) {
	db := global.GetDB()

	var users []model.User
	var total int64

	// 查询总数
	db.Model(&model.User{}).Where("`delete` = 0").Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Where("`delete` = 0").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应对象
	responses := make([]*model.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, total, nil
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(id int64) (*model.UserResponse, error) {
	db := global.GetDB()

	var user model.User
	if err := db.Where("id = ? AND `delete` = 0", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return user.ToResponse(), nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id int64, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	db := global.GetDB()
	rdb := global.GetRedis()
	ctx := context.Background()

	// 创建分布式锁，针对用户 ID 加锁
	lock := utils.NewRedisLock(rdb, fmt.Sprintf("update:user:%d", id), 10*time.Second)

	// 尝试获取锁，最多重试 3 次，每次间隔 100ms
	if err := lock.TryLock(ctx, 1, 100*time.Millisecond); err != nil {
		if errors.Is(err, utils.ErrLockFailed) {
			return nil, errors.New("系统繁忙，请稍后重试")
		}
		return nil, err
	}
	defer lock.Unlock(ctx)

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询用户是否存在
	var user model.User
	if err := tx.Where("id = ? AND `delete` = 0", id).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.NikeName != "" {
		updates["nike_name"] = req.NikeName
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	if len(updates) > 0 {
		if err := tx.Model(&user).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新查询用户
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(id int64) error {
	db := global.GetDB()

	// 查询用户是否存在
	var user model.User
	if err := db.Where("id = ? AND `delete` = 0", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 软删除
	if err := db.Model(&user).Update("delete", 1).Error; err != nil {
		return err
	}

	return nil
}
