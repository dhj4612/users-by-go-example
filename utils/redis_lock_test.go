package utils

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// 测试用的 Redis 客户端
func getTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // 使用测试数据库
	})
}

func TestRedisLock_Lock(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// 清理测试数据
	defer client.Del(ctx, "lock:test-lock")

	lock := NewRedisLock(client, "test-lock", 5*time.Second)

	// 测试获取锁
	err := lock.Lock(ctx)
	assert.NoError(t, err)

	// 测试释放锁
	err = lock.Unlock(ctx)
	assert.NoError(t, err)
}

func TestRedisLock_TryLock(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// 清理测试数据
	defer client.Del(ctx, "lock:test-trylock")

	lock := NewRedisLock(client, "test-trylock", 5*time.Second)

	// 测试带重试的获取锁
	err := lock.TryLock(ctx, 3, 100*time.Millisecond)
	assert.NoError(t, err)

	// 释放锁
	err = lock.Unlock(ctx)
	assert.NoError(t, err)
}

func TestRedisLock_Concurrent(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// 清理测试数据
	defer client.Del(ctx, "lock:test-concurrent")

	// 第一个锁获取成功
	lock1 := NewRedisLock(client, "test-concurrent", 5*time.Second)
	err := lock1.Lock(ctx)
	assert.NoError(t, err)

	// 第二个锁应该获取失败
	lock2 := NewRedisLock(client, "test-concurrent", 5*time.Second)
	err = lock2.Lock(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrLockFailed, err)

	// 释放第一个锁
	err = lock1.Unlock(ctx)
	assert.NoError(t, err)

	// 现在第二个锁应该能获取成功
	err = lock2.Lock(ctx)
	assert.NoError(t, err)

	// 释放第二个锁
	err = lock2.Unlock(ctx)
	assert.NoError(t, err)
}

func TestRedisLock_Refresh(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// 清理测试数据
	defer client.Del(ctx, "lock:test-refresh")

	lock := NewRedisLock(client, "test-refresh", 2*time.Second)

	// 获取锁
	err := lock.Lock(ctx)
	assert.NoError(t, err)

	// 等待 1 秒
	time.Sleep(1 * time.Second)

	// 续期锁
	err = lock.Refresh(ctx, 5*time.Second)
	assert.NoError(t, err)

	// 获取 TTL
	ttl, err := lock.TTL(ctx)
	assert.NoError(t, err)
	assert.True(t, ttl > 3*time.Second, "TTL should be greater than 3 seconds after refresh")

	// 释放锁
	err = lock.Unlock(ctx)
	assert.NoError(t, err)
}

func TestWithLock(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// 清理测试数据
	defer client.Del(ctx, "lock:test-withlock")

	executed := false

	// 使用 WithLock
	err := WithLock(ctx, client, "test-withlock", 5*time.Second, func() error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestWithLockRetry(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// 清理测试数据
	defer client.Del(ctx, "lock:test-withlockretry")

	executed := false

	// 使用 WithLockRetry
	err := WithLockRetry(ctx, client, "test-withlockretry", 5*time.Second, 3, 100*time.Millisecond, func() error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}
