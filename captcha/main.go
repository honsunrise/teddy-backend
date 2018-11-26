package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/config/source/file"
	"github.com/zhsyourai/teddy-backend/common/types"
	"github.com/zhsyourai/teddy-backend/common/utils"
	"google.golang.org/grpc"
	grpcHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"

	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/captcha/proto"
	"github.com/zhsyourai/teddy-backend/captcha/repositories"
	"github.com/zhsyourai/teddy-backend/captcha/server"
	"github.com/zhsyourai/teddy-backend/common/config"
)

func main() {
	conf, err := config.NewConfig(file.NewSource(file.WithFormat(config.Yaml), file.WithPath("config/config.yaml")))
	if err != nil {
		log.Fatal(err)
	}
	var confType types.Config
	err = conf.Scan(&confType)
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
	kvRepo, err := repositories.NewKeyValuePairRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}

	// New Handler
	captchaSrv, err := server.NewCaptchaServer(kvRepo)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterCaptchaServer(grpcServer, captchaSrv)

	healthSrv := grpcHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// Run service
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
