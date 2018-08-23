package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	"github.com/zhsyourai/teddy-backend/uaa/client"
	"github.com/zhsyourai/teddy-backend/uaa/handler/account"
	uaa "github.com/zhsyourai/teddy-backend/uaa/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.srv.uaa"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the Example srv client
		micro.WrapHandler(client.NotifyWrapper(service)))

	// Register Handler
	uaa.RegisterUAAHandler(service.Server(), account.GetInstance())

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
