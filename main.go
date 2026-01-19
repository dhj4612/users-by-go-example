package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	app "users-by-go-example/application"
	"users-by-go-example/internal/router"
)

func main() {
	// 初始化外部资源
	app.InitAll()

	// 延迟关闭资源
	defer app.Close()

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	conf := app.GetConfig()
	addr := fmt.Sprintf(":%d", conf.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		fmt.Printf("服务器启动在端口 %d\n", conf.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 设置 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已退出")
}
