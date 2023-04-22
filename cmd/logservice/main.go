package main

import (
	"context"
	"fmt"

	stlog "log"

	"github.com/Tacks9/go-distributed/log"
	"github.com/Tacks9/go-distributed/registry"
	"github.com/Tacks9/go-distributed/service"
)

func main() {
	// 申请一个日志地址
	log.Run("./go-distributed.log")

	// 日志服务的地址
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	// 日志服务
	regItem := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       serviceAddress,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	// 启动 Log 服务
	ctx, err := service.Start(context.Background(),
		host,
		port,
		regItem,
		log.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("Shutting down log service.")
}
