package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Tacks9/go-distributed/registry"
)

func main() {

	// 启动注册服务-服务端
	http.Handle("/services", &registry.RegistrationService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServicePort

	// 注册服务-遇到异常停止
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	// 手动关闭
	go func() {
		fmt.Println("Registry Service Started. Press any Key to Stop.")
		var s string
		fmt.Scanln(&s)
		// 关闭服务
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()

	fmt.Println("Shutting Down Registry Service...")
}
