package main

import (
	"context"
	"fmt"

	stlog "log"

	"github.com/Tacks9/go-distributed/grades"
	"github.com/Tacks9/go-distributed/log"
	"github.com/Tacks9/go-distributed/registry"
	"github.com/Tacks9/go-distributed/service"
)

func main() {
	// grading 服务的地址
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	// 注册服务
	regItem := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceURL:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
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

	// 获取到日志服务
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s \n", logProvider)
		log.SetClientLogger(logProvider, regItem.ServiceName)
		// 申请一个日志地址
	}
	<-ctx.Done()

	fmt.Println("Shutting down grading service.")
}
