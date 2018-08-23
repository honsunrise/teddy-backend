package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/zhsyourai/teddy-backend/notify/handler/notify"
	"github.com/zhsyourai/teddy-backend/notify/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.srv.notify"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	proto.RegisterNotifyHandler(service.Server(), notify.GetInstance())

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
