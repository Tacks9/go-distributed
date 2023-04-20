package main

import (
	"context"
	"fmt"

	stlog "log"

	"github.com/Tacks9/go-distributed/grades"
	"github.com/Tacks9/go-distributed/registry"
	"github.com/Tacks9/go-distributed/service"
)

func main() {
	// grading 服务的地址
	host, port := "localhost", "6000"

	// 注册服务
	regItem := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
	}
	// 启动服务
	ctx, err := service.Start(context.Background(),
		host,
		port,
		regItem,
		grades.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("Shutting down grading service.")
}
