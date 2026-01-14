# Redis 分布式锁使用指南

本项目使用 `bsm/redislock` 库实现 Redis 分布式锁。

## 基本用法

### 1. 标准用法（手动管理锁）

```go
package main

import (
    "context"
    "time"
    "users-by-go-example/global"
    "users-by-go-example/utils"
)

func example1() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    // 创建锁
    lock := utils.NewRedisLock(rdb, "my-resource", 10*time.Second)

    // 获取锁
    if err := lock.Lock(ctx); err != nil {
        if errors.Is(err, utils.ErrLockFailed) {
            return errors.New("资源正在被使用")
        }
        return err
    }
    defer lock.Unlock(ctx)

    // 执行业务逻辑
    // ...

    return nil
}
```

### 2. 带重试的用法

```go
func example2() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    lock := utils.NewRedisLock(rdb, "my-resource", 10*time.Second)

    // 尝试获取锁，最多重试 3 次，每次间隔 100ms
    if err := lock.TryLock(ctx, 3, 100*time.Millisecond); err != nil {
        return err
    }
    defer lock.Unlock(ctx)

    // 执行业务逻辑
    // ...

    return nil
}
```

### 3. 便捷函数用法（推荐）

```go
func example3() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    // 使用 WithLock 自动管理锁的获取和释放
    return utils.WithLock(ctx, rdb, "my-resource", 10*time.Second, func() error {
        // 执行业务逻辑
        // ...
        return nil
    })
}

func example4() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    // 使用 WithLockRetry 支持重试
    return utils.WithLockRetry(
        ctx,
        rdb,
        "my-resource",
        10*time.Second,
        3,                      // 重试 3 次
        100*time.Millisecond,   // 每次间隔 100ms
        func() error {
            // 执行业务逻辑
            // ...
            return nil
        },
    )
}
```

## 高级功能

### 1. 锁续期

```go
func exampleRefresh() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    lock := utils.NewRedisLock(rdb, "long-task", 10*time.Second)
    if err := lock.Lock(ctx); err != nil {
        return err
    }
    defer lock.Unlock(ctx)

    // 执行一些操作
    doSomething()

    // 延长锁的过期时间
    if err := lock.Refresh(ctx, 10*time.Second); err != nil {
        return err
    }

    // 继续执行
    doMoreWork()

    return nil
}
```

### 2. 查询锁的剩余时间

```go
func exampleTTL() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    lock := utils.NewRedisLock(rdb, "my-resource", 10*time.Second)
    if err := lock.Lock(ctx); err != nil {
        return err
    }
    defer lock.Unlock(ctx)

    // 获取锁的剩余时间
    ttl, err := lock.TTL(ctx)
    if err != nil {
        return err
    }

    fmt.Printf("锁剩余时间: %v\n", ttl)

    return nil
}
```

### 3. 获取锁的唯一标识

```go
func exampleMetadata() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    lock := utils.NewRedisLock(rdb, "my-resource", 10*time.Second)
    if err := lock.Lock(ctx); err != nil {
        return err
    }
    defer lock.Unlock(ctx)

    // 获取锁的唯一标识（用于调试）
    metadata := lock.Metadata()
    fmt.Printf("锁标识: %s\n", metadata)

    return nil
}
```

## 实际应用场景

### 场景 1：防止重复注册

```go
func (s *UserService) Register(req *model.RegisterRequest) (*model.UserResponse, error) {
    ctx := context.Background()
    rdb := global.GetRedis()

    // 针对用户名加锁，防止并发注册
    lock := utils.NewRedisLock(rdb, "register:"+req.Username, 10*time.Second)
    if err := lock.TryLock(ctx, 3, 100*time.Millisecond); err != nil {
        if errors.Is(err, utils.ErrLockFailed) {
            return nil, errors.New("系统繁忙，请稍后重试")
        }
        return nil, err
    }
    defer lock.Unlock(ctx)

    // 检查用户名是否存在
    // 创建用户
    // ...
}
```

### 场景 2：防止并发更新

```go
func (s *UserService) UpdateUser(id int64, req *model.UpdateUserRequest) (*model.UserResponse, error) {
    ctx := context.Background()
    rdb := global.GetRedis()

    // 针对用户 ID 加锁，防止并发更新
    lockKey := fmt.Sprintf("update:user:%d", id)

    return utils.WithLockRetry(
        ctx,
        rdb,
        lockKey,
        10*time.Second,
        3,
        100*time.Millisecond,
        func() error {
            // 查询用户
            // 更新用户
            // ...
            return nil
        },
    )
}
```

### 场景 3：定时任务防止重复执行

```go
func CronJob() error {
    ctx := context.Background()
    rdb := global.GetRedis()

    // 使用锁确保定时任务只在一个实例上执行
    lock := utils.NewRedisLock(rdb, "cron:daily-report", 5*time.Minute)
    if err := lock.Lock(ctx); err != nil {
        if errors.Is(err, utils.ErrLockFailed) {
            // 其他实例正在执行，直接返回
            return nil
        }
        return err
    }
    defer lock.Unlock(ctx)

    // 执行定时任务
    generateDailyReport()

    return nil
}
```

### 场景 4：库存扣减

```go
func DeductStock(productID int64, quantity int) error {
    ctx := context.Background()
    rdb := global.GetRedis()

    lockKey := fmt.Sprintf("stock:%d", productID)

    return utils.WithLockRetry(
        ctx,
        rdb,
        lockKey,
        5*time.Second,
        5,                      // 库存操作重试 5 次
        50*time.Millisecond,
        func() error {
            // 查询库存
            stock := getStock(productID)

            // 检查库存是否足够
            if stock < quantity {
                return errors.New("库存不足")
            }

            // 扣减库存
            updateStock(productID, stock-quantity)

            return nil
        },
    )
}
```

## 最佳实践

### 1. 锁的粒度

```go
// ❌ 不好：锁粒度太大
lock := utils.NewRedisLock(rdb, "all-users", 10*time.Second)

// ✅ 好：针对具体资源加锁
lock := utils.NewRedisLock(rdb, "user:"+username, 10*time.Second)
```

### 2. 锁的过期时间

```go
// 根据业务场景设置合理的过期时间
// 快速操作：1-5 秒
lock := utils.NewRedisLock(rdb, "quick-op", 3*time.Second)

// 普通操作：5-30 秒
lock := utils.NewRedisLock(rdb, "normal-op", 10*time.Second)

// 长时间操作：30 秒以上，建议使用 Refresh
lock := utils.NewRedisLock(rdb, "long-op", 30*time.Second)
```

### 3. 错误处理

```go
lock := utils.NewRedisLock(rdb, "my-resource", 10*time.Second)
if err := lock.Lock(ctx); err != nil {
    if errors.Is(err, utils.ErrLockFailed) {
        // 锁被占用，返回友好提示
        return errors.New("系统繁忙，请稍后重试")
    }
    // 其他错误（Redis 连接失败等）
    return fmt.Errorf("获取锁失败: %w", err)
}
defer lock.Unlock(ctx)
```

### 4. 使用 defer 确保锁释放

```go
// ✅ 好：使用 defer 确保锁一定会被释放
lock := utils.NewRedisLock(rdb, "my-resource", 10*time.Second)
if err := lock.Lock(ctx); err != nil {
    return err
}
defer lock.Unlock(ctx)

// 即使发生 panic，defer 也会执行
doSomethingRisky()
```

## 注意事项

1. **锁的过期时间要大于业务执行时间**，否则可能导致锁提前释放
2. **使用 defer 释放锁**，确保即使发生错误也能释放
3. **锁的 key 要有业务含义**，方便调试和监控
4. **避免死锁**：不要在持有锁的情况下再次获取同一个锁
5. **合理设置重试次数**：避免无限重试导致请求堆积

## 与手动实现的对比

| 特性 | 手动实现 | bsm/redislock |
|------|---------|---------------|
| 唯一标识 | time.Now().String() | UUID（更可靠） |
| 锁续期 | ❌ | ✅ |
| TTL 查询 | ❌ | ✅ |
| 重试策略 | 简单 | 灵活（线性/指数退避） |
| 代码维护 | 需要自己维护 | 库维护 |
| 生产就绪 | 需要改进 | ✅ |

## 相关资源

- [bsm/redislock GitHub](https://github.com/bsm/redislock)
- [Redis 分布式锁最佳实践](https://redis.io/docs/manual/patterns/distributed-locks/)
