package utils

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrLockFailed   = errors.New("获取锁失败")
	ErrUnlockFailed = errors.New("释放锁失败")
)

// RedisLock Redis 分布式锁
type RedisLock struct {
	client *redis.Client
	key    string
	value  string
	expire time.Duration
}

// NewRedisLock 创建 Redis 分布式锁
func NewRedisLock(client *redis.Client, key string, expire time.Duration) *RedisLock {
	return &RedisLock{
		client: client,
		key:    "lock:" + key,
		value:  generateLockValue(),
		expire: expire,
	}
}

// generateLockValue 生成锁的唯一值
func generateLockValue() string {
	return time.Now().String()
}

// Lock 获取锁
func (l *RedisLock) Lock(ctx context.Context) error {
	// 使用 SET NX EX 命令，只有当 key 不存在时才设置，并设置过期时间
	success, err := l.client.SetNX(ctx, l.key, l.value, l.expire).Result()
	if err != nil {
		return err
	}
	if !success {
		return ErrLockFailed
	}
	return nil
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	// 使用 Lua 脚本确保只删除自己持有的锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}
	if result == int64(0) {
		return ErrUnlockFailed
	}
	return nil
}

// TryLock 尝试获取锁，支持重试
func (l *RedisLock) TryLock(ctx context.Context, retryTimes int, retryDelay time.Duration) error {
	for i := 0; i < retryTimes; i++ {
		err := l.Lock(ctx)
		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrLockFailed) {
			return err
		}
		// 等待后重试
		time.Sleep(retryDelay)
	}
	return ErrLockFailed
}
