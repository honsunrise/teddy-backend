package main

import (
	"github.com/casbin/casbin"
	"github.com/casbin/mongodb-adapter"
	"github.com/micro/go-micro"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/utils"
	"github.com/zhsyourai/teddy-backend/uaa/components"

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
	// Load config
	conf := config.GetConfig()
	mongodbUri := utils.BuildMongodbURI(conf.Databases["mongodb"])
	model := casbin.NewModel(conf.Casbin)

	enforcer := casbin.NewEnforcer(model, mongodbadapter.NewAdapter(mongodbUri))
	enforcer.LoadPolicy()

	// New Mongodb client
	mongodbClient, err := mongo.Connect(context.Background(), mongodbUri)
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
	// New components
	uidGenerator, err := components.NewUidGenerator(accountRepo)
	if err != nil {
		log.Fatal(err)
	}
	// New Handler
	accountHandler, err := account.NewAccountHandler(accountRepo, uidGenerator, enforcer)
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
