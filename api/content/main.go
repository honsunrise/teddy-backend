package main

import (
	"github.com/micro/go-log"

	"github.com/micro/go-micro"
	"github.com/zhsyourai/teddy-backend/api/content/client"
	"github.com/zhsyourai/teddy-backend/api/content/handler"

	"github.com/zhsyourai/teddy-backend/api/content/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.content"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the Content srv client
		micro.WrapHandler(client.MessageWrapper(service)),
	)

	// Register Handler
	proto.RegisterContentHandler(service.Server(), new(handler.Content))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
