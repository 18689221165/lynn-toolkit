package web

import (
	"context"
	"fmt"
	"github.com/18689221165/lynn-toolkit/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config web server config
type Config struct {
	RunMode         string `yaml:"runMode"`         // 启动模式:debug|release|test
	Port            int    `yaml:"port"`            // 服务器端口
	ShutdownTimeout int    `yaml:"shutdownTimeout"` // 优雅停止服务的超时时间（秒）
}

type GinServer struct {
	conf   Config
	engine *gin.Engine
	server *http.Server
	logger *zap.SugaredLogger
}

func NewGinServer(conf Config, log *zap.SugaredLogger) *GinServer {

	gin.DisableConsoleColor()
	gin.SetMode(conf.RunMode)

	engine := gin.New()
	engine.Use(IgnoreIndexAndFavicon(), GinZapLog(log.Desugar(), conf.RunMode), RecoveryWithZap(log.Desugar(), true))

	// 空路径响应
	engine.NoRoute(func(c *gin.Context) { c.JSON(http.StatusOK, service.ErrNotFound) })
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", conf.Port), Handler: engine}

	return &GinServer{conf: conf, engine: engine, server: httpServer, logger: log}
}

// RouteGroup 增加组路由
func (instance *GinServer) RouteGroup(relativePath string, middlewares ...gin.HandlerFunc) *gin.RouterGroup {
	return instance.engine.Group(relativePath, middlewares...)
}

// Static 增加静态文件
func (instance *GinServer) Static(relativePath, root string) {
	instance.engine.Static(relativePath, root)
}

// RunAsync 异步方式启动
func (instance *GinServer) RunAsync() {
	go func() {
		err := instance.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			instance.logger.Fatalf("Web 服务器实例启动失败，原因：%+v", err)
		}
	}()
	instance.logger.Infof("Web 服务器实例启动成功，访问端口[%d]，耶！", instance.conf.Port)
}

// WaitInterrupt 等待接受中断信号
// 1、停止接收新请求，等待已有请求执行完毕；
// 2、如果等待时间超过ShutdownTimeout设定的时间，则强制关闭HttpServer
func (instance *GinServer) WaitInterrupt() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}

// Shutdown 优雅关闭WebServer
func (instance *GinServer) Shutdown() {
	instance.logger.Warnf("Web 服务器实例正在停止,倒计时[%d]秒...", instance.conf.ShutdownTimeout)

	// HttpServer退出等待时间，默认等待10秒
	timeout := 10 * time.Second
	if instance.conf.ShutdownTimeout > 0 {
		timeout = time.Duration(instance.conf.ShutdownTimeout) * time.Second
	}
	// 设置退出等待超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// 关闭HttpServer
	if err := instance.server.Shutdown(ctx); err != nil {
		instance.logger.Warnf("Web 服务器实例无法停止, 错误原因： %+v", err)
		return
	}
	instance.logger.Warnf("Web 服务器实例已停止，88！")
}

func (instance *GinServer) AddMiddleware(middleware ...gin.HandlerFunc) {
	instance.engine.Use(middleware...)
}

func (instance GinServer) LoadHTMLGlob(path string) {
	instance.engine.LoadHTMLGlob(path)
}

func (instance GinServer) LoadHTMLFiles(files ...string) {
	instance.engine.LoadHTMLFiles(files...)
}
