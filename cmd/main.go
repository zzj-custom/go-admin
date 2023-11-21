package main

import (
	_ "admin/internal/controller"
	"admin/internal/initializer"
	"context"
	"log/slog"
	"os"
)

func main() {
	// 初始化日志
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(textLogger)

	// 初始化配置文件
	initializer.Config()

	// 假如上下文取消处理
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	defer func() {
		cancelFunc()
	}()

	// 初始化资源
	initializer.InitResource()

	initializer.Server(cancelCtx)
	slog.Info("所有服务已经关闭")
}
