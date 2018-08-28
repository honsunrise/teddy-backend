package main

import (
	"context"
	"flag"
	"github.com/micro/go-micro"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/utils"
	"github.com/zhsyourai/teddy-backend/content/handler/content"
	"github.com/zhsyourai/teddy-backend/content/proto"
	"github.com/zhsyourai/teddy-backend/content/repositories"
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
	contentRepo, err := repositories.NewContentRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.srv.content"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	// New Handler
	contentHandler, err := content.NewContentHandler(contentRepo)
	if err != nil {
		log.Fatal(err)
	}
	// Register Handler
	proto.RegisterContentHandler(service.Server(), contentHandler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
