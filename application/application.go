package application

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
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
	instance  = content{}
	once      sync.Once
	closeOnce sync.Once
)

// InitAll 初始化所有资源
func InitAll() {
	once.Do(func() {
		initConfig()
		initDB()
		initRedis()
	})
}

// InitConfig 初始化配置
func initConfig() {
	conf := config.Load()
	instance.Config = conf
}

// InitDB 初始化数据库连接
func initDB() {
	conf := instance.Config
	connectUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(connectUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)                 // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                // 最大打开连接数
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接空闲最大存活时间

	instance.DB = db
	log.Println("数据库连接成功")
}

// InitRedis 初始化 Redis 连接
func initRedis() {
	conf := instance.Config
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
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

// GetConfig 获取配置
func GetConfig() *config.Config {
	if instance.Config == nil {
		panic("配置未初始化")
	}
	return instance.Config
}

var apiPermitsMap map[string]string = nil
var getApiPermitsMapOnce sync.Once

// GetApiPermitsMap 获取Api权限标识Map
func GetApiPermitsMap() map[string]string {
	if apiPermitsMap != nil {
		return apiPermitsMap
	}

	if instance.Config == nil {
		panic("配置未初始化")
	}

	getApiPermitsMapOnce.Do(func() {
		apiPermitsMap = make(map[string]string)
		for _, item := range GetConfig().ApiPermits {
			apiPermitsMap[item.Method+" "+item.Path] = item.Permits
		}
	})

	return apiPermitsMap
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	if instance.DB == nil {
		panic("数据库未初始化")
	}
	return instance.DB
}

// GetRedis 获取 Redis 连接
func GetRedis() *redis.Client {
	if instance.Redis == nil {
		panic("Redis 未初始化")
	}
	return instance.Redis
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

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if instance.Redis != nil {
		return instance.Redis.Close()
	}
	return nil
}

// Close 关闭所有资源
func Close() {
	closeOnce.Do(func() {
		if err := CloseDB(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		}
		if err := CloseRedis(); err != nil {
			log.Printf("关闭 Redis 连接失败: %v", err)
		}

		instance.DB = nil
		instance.Redis = nil
	})
}
