package app

import (
	"context"
	"github.com/guidoxie/knife/pkg/app/http"
	"github.com/guidoxie/knife/pkg/log"
	stdHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type IApp interface {
	// 获取应用程序名称。
	Name() string
	// 获取应用程序 Http Server 监听地址。
	HttpAddr() string
	// 获取 IHttpServer 对象，用于注册路由，设置中间件等
	HttpServer() http.IServer
	// 启动
	Run()
	// 停止
	Stop(ctx context.Context) error
}

var app *App

func Current() *App {
	return app
}

type App struct {
	name       string
	httpServer *http.Server
}

func (a *App) HttpServer() http.IServer {
	return a.httpServer
}

func New(name string, options ...http.Options) *App {
	app = &App{
		name:       name,
		httpServer: http.NewServer(name, options...),
	}
	// 添加默认中间件
	app.httpServer.AddDefaultMiddle(
		http.AccessLog(),
		http.Recovery(),
		http.Translations())
	return app
}

func (a *App) Name() string {
	return a.name
}

func (a *App) HttpAddr() string {
	return a.httpServer.Addr
}

func (a *App) Run() {
	go func() {
		if err := a.httpServer.Run(); err != nil && err != stdHttp.ErrServerClosed {
			log.Fatal("HTTP server err:", err)
		}
	}()
	// 监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Info("Shutdown server...")
	if err := a.Stop(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Info("Server exiting")
	a.AfterStop()
}

func (a *App) Stop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return a.httpServer.Shutdown(ctx)
}

func (a *App) AfterStop() {
	log.Sync()
}
