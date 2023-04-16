package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Tacks9/go-distributed/registry"
)

// 服务注册
func Start(ctx context.Context, host, port string,
	reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {

	// 设置路由和Handler
	registerHandlersFunc()

	// 启动服务：启动HTTP服务返回新的上下文
	ctx = startService(ctx, reg.ServiceName, host, port)

	// 注册服务：将服务信息注册到注册中心
	err := registry.RegistryService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil

}

// 服务启动
func startService(parent context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	// 创建一个具有取消功能的 上下文
	ctx, cancel := context.WithCancel(parent)

	var srv http.Server
	srv.Addr = ":" + port

	fmt.Printf("HTTP Server listening on port %s...\n", port)

	// 启动 HTTP 服务
	go func() {
		// 启动 HTTP 服务 如果发生错误，则 cancel()
		log.Println(srv.ListenAndServe())
		// 如果异常 执行服务移除
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		// 发送信号 取消上下文
		cancel()
	}()

	// 用户终止
	go func() {
		fmt.Printf("[%v] started. \nPress any key to stop. \n", serviceName)
		// 等待用户输入
		var s string
		fmt.Scanln(&s)

		// 如果异常 执行服务移除
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		// 手动停止 发送信号
		srv.Shutdown(ctx)
		// 取消上下文
		cancel()
	}()

	return ctx
}
