package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/zhsyourai/teddy-backend/common/utils"

	"context"
	"flag"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/uaa/client"
	"github.com/zhsyourai/teddy-backend/uaa/handler/account"
	uaa "github.com/zhsyourai/teddy-backend/uaa/proto"
	"github.com/zhsyourai/teddy-backend/uaa/repositories"
)

func main() {
	flag.Parse()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	// New Mongodb client
	conf := config.GetConfig()
	mongodbClient, err := mongo.Connect(context.Background(), utils.BuildMongodbURI(conf.Databases["mongodb"]))
	if err != nil {
		log.Fatal(err)
	}
	// New Repository
	accountRepo, err := repositories.NewAccountRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}
	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.srv.uaa"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the Example srv client
		micro.WrapHandler(client.NotifyWrapper(service)))
	// New Handler
	accountHandler, err := account.NewAccountHandler(accountRepo)
	if err != nil {
		log.Fatal(err)
	}
	// Register Handler
	uaa.RegisterUAAHandler(service.Server(), accountHandler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
