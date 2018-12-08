package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/config/source/file"
	"github.com/zhsyourai/teddy-backend/common/utils"
	"github.com/zhsyourai/teddy-backend/uaa/components"
	"google.golang.org/grpc"
	grpcHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"

	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
	"github.com/zhsyourai/teddy-backend/uaa/repositories"
	"github.com/zhsyourai/teddy-backend/uaa/server"
)

func init() {
	log.SetReportCaller(true)
}

func main() {
	conf, err := config.NewConfig(file.NewSource(file.WithFormat(config.Yaml), file.WithPath("config/config.yaml")))
	if err != nil {
		log.Fatal(err)
	}
	confSecret, err := config.NewConfig(file.NewSource(file.WithFormat(config.Yaml), file.WithPath("secret/config.yaml")))
	if err != nil {
		log.Fatal(err)
	}
	var confType Config
	err = conf.Scan(&confType)
	if err != nil {
		log.Fatal(err)
	}
	err = confSecret.Scan(&confType)
	if err != nil {
		log.Fatal(err)
	}
	// Load config
	mongodbUri := utils.BuildMongodbURI(confType.Databases["mongodb"])

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

	// New components
	uidGenerator, err := components.NewUidGenerator(accountRepo)
	if err != nil {
		log.Fatal(err)
	}
	// New Handler
	accountSrv, err := server.NewAccountServer(accountRepo, uidGenerator)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUAAServer(grpcServer, accountSrv)

	healthSrv := grpcHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// Run service
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
