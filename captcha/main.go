package main

import (
	"context"
	"flag"
	"github.com/micro/go-micro"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/captcha/handler/captcha"
	"github.com/zhsyourai/teddy-backend/captcha/proto"
	"github.com/zhsyourai/teddy-backend/captcha/repositories"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/utils"
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
	contentRepo, err := repositories.NewKeyValuePairRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.srv.captcha"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	// New Handler
	contentHandler, err := captcha.NewCaptchaHandler(contentRepo)
	if err != nil {
		log.Fatal(err)
	}
	// Register Handler
	proto.RegisterCaptchaHandler(service.Server(), contentHandler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
