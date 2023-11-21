package initializer

import (
	"admin/cmd/middleware"
	"admin/cmd/router"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func Server(ctx context.Context) {
	wg.Add(1)
	// 初始化gin
	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()

	// 设置全局中间件
	e.Use(middleware.DefaultLimit)

	// 初始化路由
	e = router.InitRouter(e)

	// 启动监听
	srv := &http.Server{
		Addr:              viper.GetString("server.host"),
		Handler:           e,
		ReadTimeout:       viper.GetDuration("server.read_timeout"),
		ReadHeaderTimeout: viper.GetDuration("server.read_header_timeout"),
		WriteTimeout:      viper.GetDuration("write_timeout"),
	}
	srv.SetKeepAlivesEnabled(true)

	go func() {
		defer func() {
			wg.Done()
		}()
		// 启动信号监听实现程序阻塞
		// 信号：0x02(INT), 0x09(KILL), 0x0F(TERM)
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt, os.Kill, os.Signal(syscall.SIGTERM))
		sig := <-quit

		// 收到监听的信号，准备退出程序
		log.Printf("收到退出程序信号: %d", sig)

		// 加入超时上下文处理程序退出
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("关闭服务失败，错误：%v", err)
			return
		}
		log.Println("Web服务已关闭")
	}()

	// 启用协程运行服务监听
	go func() {
		log.Printf("启动服务，监听地址：%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("启动服务失败或服务运行时出现异常，错误: %v", err)
		}
	}()

	wg.Wait()
}
