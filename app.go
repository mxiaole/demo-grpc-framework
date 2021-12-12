package main

import (
	"fmt"
	"google.golang.org/grpc"
	"kratos-demo/proto"
	"kratos-demo/server"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// 定义应用

type App struct {
	// 必须设置的值
	Id      string
	Name    string // 服务名
	Version string // 版本号

	// 可选的值
	opts Option // 其他参数
}

// Option 定义可选参数有哪些
type Option struct {
	register string
	balance  string
}

type option func(*Option)

func Register(register string) option {
	return func(opt *Option) {
		opt.register = register
	}
}
func Balance(balance string) option {
	return func(opt *Option) {
		opt.balance = balance
	}
}

// 一个应用应该有哪些方法

// Start 启动应用
func (a *App) Start() {
	// 启动一个人grpc的服务
	lis, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("listen error.......", err)
	}
	// 创建服务
	grpcServer := grpc.NewServer()

	// 给grpc服务添加路由
	proto.RegisterGreeterServer(grpcServer, &server.GreetServe{})

	// 开启一个协程进行服务监听
	go func() error {
		log.Println("server start.........")
		return grpcServer.Serve(lis)
	}()

	// 1. 如何将服务注册到注册中心呢？
	// 什么时候将服务的信息上报到注册中心呢？要在服务正常启动之后才可以上报
	if a.opts.register != "" {
		fmt.Printf("将服务: %s, 注册到: %s\n", a.Name, a.opts.register)
	}

	// 使用信号来停止服务
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT) // 使用15号信号，该信号可以被捕获
	done := make(chan bool, 1)
	go func() {
		s := <-c
		fmt.Printf("get signal %v\n", s)
		done <- true
	}()

	<-done
	a.Stop()
}

// Stop 优雅关闭服务
func (a *App) Stop() {
	// 主要用来做一些资源清理的工作
	fmt.Println("kill signal killed the server")
}

// NewApp 创建一个应用
func NewApp(id, name, version string, option ...option) *App {
	// 给app设置配置选项
	opts := Option{}
	for _, o := range option {
		o(&opts)
	}

	return &App{
		Id:      id,
		Name:    name,
		Version: version,
		opts:    opts,
	}
}

func main() {
	options := []option{Register("etd"), Balance("p2c")}
	app := NewApp("1", "demo", "v1.1.0", options...)
	fmt.Println(app)
	app.Start()
}
