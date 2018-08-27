package main

import (
	"context"
	"flag"
	"github.com/micro/go-micro"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/utils"
	"github.com/zhsyourai/teddy-backend/notify/handler/notify"
	"github.com/zhsyourai/teddy-backend/notify/proto"
	"github.com/zhsyourai/teddy-backend/notify/repositories"
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
	notifyRepo, err := repositories.NewInBoxRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.srv.notify"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	// New Handler
	notifyHandler, err := notify.NewNotifyHandler(notifyRepo)
	if err != nil {
		log.Fatal(err)
	}
	// Register Handler
	proto.RegisterNotifyHandler(service.Server(), notifyHandler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
