package main

import (
	"github.com/micro/go-log"

	"github.com/micro/go-micro"
	"github.com/zhsyourai/teddy-backend/api/message/client"
	"github.com/zhsyourai/teddy-backend/api/message/handler"

	"github.com/zhsyourai/teddy-backend/api/message/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.message"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the Message srv client
		micro.WrapHandler(client.MessageWrapper(service)),
	)

	// Register Handler
	proto.RegisterMessageHandler(service.Server(), new(handler.Message))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
