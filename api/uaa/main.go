package main

import (
	"github.com/micro/go-log"

	"github.com/micro/go-micro"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/api/uaa/handler"

	"github.com/zhsyourai/teddy-backend/api/uaa/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.api.uaa"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the UaaApi srv client
		micro.WrapHandler(client.UaaWrapper(service)),
	)

	// Register Handler
	proto.RegisterUaaApiHandler(service.Server(), new(handler.UaaApi))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
