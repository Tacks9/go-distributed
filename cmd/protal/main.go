package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/Tacks9/go-distributed/log"
	"github.com/Tacks9/go-distributed/portal"
	"github.com/Tacks9/go-distributed/registry"
	"github.com/Tacks9/go-distributed/service"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}

	host, port := "localhost", "8080"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	regItem := registry.Registration{
		ServiceName: registry.ProtalService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.GradingService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(),
		host,
		port,
		regItem,
		portal.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	// 获取日志服务
	if logProvider, err := registry.GetProvider(registry.LogService); err != nil {
		log.SetClientLogger(logProvider, regItem.ServiceName)
	}

	<-ctx.Done()

	fmt.Println("Shutting down portal.")

}
