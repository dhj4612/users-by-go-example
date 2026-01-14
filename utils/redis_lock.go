package utils

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

var (
	ErrLockFailed   = errors.New("获取锁失败")
	ErrUnlockFailed = errors.New("释放锁失败")
)

// RedisLock Redis 分布式锁（基于 bsm/redislock）
type RedisLock struct {
	locker *redislock.Client
	lock   *redislock.Lock
	key    string
	expire time.Duration
}

// NewRedisLock 创建 Redis 分布式锁
func NewRedisLock(client *redis.Client, key string, expire time.Duration) *RedisLock {
	return &RedisLock{
		locker: redislock.New(client),
		key:    "lock:" + key,
		expire: expire,
	}
}

// Lock 获取锁
func (l *RedisLock) Lock(ctx context.Context) error {
	lock, err := l.locker.Obtain(ctx, l.key, l.expire, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return ErrLockFailed
	}
	if err != nil {
		return err
	}
	l.lock = lock
	return nil
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	if l.lock == nil {
		return ErrUnlockFailed
	}
	if err := l.lock.Release(ctx); err != nil {
		return err
	}
	return nil
}

// TryLock 尝试获取锁，支持重试
func (l *RedisLock) TryLock(ctx context.Context, retryTimes int, retryDelay time.Duration) error {
	// 使用 bsm/redislock 的重试策略
	lock, err := l.locker.Obtain(ctx, l.key, l.expire, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(
			redislock.LinearBackoff(retryDelay),
			retryTimes,
		),
	})
	if errors.Is(err, redislock.ErrNotObtained) {
		return ErrLockFailed
	}
	if err != nil {
		return err
	}
	l.lock = lock
	return nil
}

// Refresh 延长锁的过期时间
func (l *RedisLock) Refresh(ctx context.Context, expire time.Duration) error {
	if l.lock == nil {
		return errors.New("锁未持有")
	}
	if err := l.lock.Refresh(ctx, expire, nil); err != nil {
		return err
	}
	return nil
}

// TTL 获取锁的剩余时间
func (l *RedisLock) TTL(ctx context.Context) (time.Duration, error) {
	if l.lock == nil {
		return 0, errors.New("锁未持有")
	}
	ttl, err := l.lock.TTL(ctx)
	if err != nil {
		return 0, err
	}
	return ttl, nil
}

// Metadata 获取锁的元数据（唯一标识）
func (l *RedisLock) Metadata() string {
	if l.lock == nil {
		return ""
	}
	return l.lock.Metadata()
}

// ============ 便捷函数 ============

// WithLock 使用锁执行函数（自动获取和释放锁）
func WithLock(ctx context.Context, client *redis.Client, key string, expire time.Duration, fn func() error) error {
	lock := NewRedisLock(client, key, expire)
	if err := lock.Lock(ctx); err != nil {
		return err
	}
	defer lock.Unlock(ctx)
	return fn()
}

// WithLockRetry 使用锁执行函数（支持重试）
func WithLockRetry(ctx context.Context, client *redis.Client, key string, expire time.Duration, retryTimes int, retryDelay time.Duration, fn func() error) error {
	lock := NewRedisLock(client, key, expire)
	if err := lock.TryLock(ctx, retryTimes, retryDelay); err != nil {
		return err
	}
	defer lock.Unlock(ctx)
	return fn()
}
