package global

import (
	"context"
	"fmt"
	"log"
	"sync"
	"users-by-go-example/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// content 全局资源管理
type content struct {
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

var (
	instance content
	once     sync.Once
)

// InitConfig 初始化配置
func InitConfig() {
	cfg := config.Load()
	instance.Config = cfg
}

// GetConfig 获取配置
func GetConfig() *config.Config {
	return instance.Config
}

// InitDB 初始化数据库连接
func InitDB() {
	cfg := instance.Config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	instance.DB = db
	log.Println("数据库连接成功")
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return instance.DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if instance.DB != nil {
		sqlDB, err := instance.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// InitRedis 初始化 Redis 连接
func InitRedis() {
	cfg := instance.Config
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	instance.Redis = rdb
	log.Println("Redis 连接成功")
}

// GetRedis 获取 Redis 连接
func GetRedis() *redis.Client {
	return instance.Redis
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if instance.Redis != nil {
		return instance.Redis.Close()
	}
	return nil
}

// Close 关闭所有资源
func Close() {
	if err := CloseDB(); err != nil {
		log.Printf("关闭数据库连接失败: %v", err)
	}
	if err := CloseRedis(); err != nil {
		log.Printf("关闭 Redis 连接失败: %v", err)
	}
}
