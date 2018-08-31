package main

import (
	"github.com/casbin/casbin"
	"github.com/casbin/mongodb-adapter"
	"github.com/micro/go-log"
	"github.com/zhsyourai/teddy-backend/api/uaa/handler"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/utils"

	"github.com/micro/go-micro"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/api/uaa/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.uaa"),
		micro.Version("latest"),
	)

	// Load config
	conf := config.GetConfig()
	mongodbUri := utils.BuildMongodbURI(conf.Databases["mongodb"])
	model := casbin.NewModel(conf.Casbin)

	enforcer := casbin.NewEnforcer(model, mongodbadapter.NewAdapter(mongodbUri))
	enforcer.LoadPolicy()

	// Initialise service
	service.Init(
		// create wrap for the Message srv client
		micro.WrapHandler(client.UaaWrapper(service)),
	)
	uaa, err := handler.NewUaaHandler(enforcer)
	if err != nil {
		log.Fatal(err)
	}

	// Register Handler
	proto.RegisterUaaHandler(service.Server(), uaa)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
