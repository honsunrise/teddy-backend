package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	grpcHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	captchaProto "teddy-backend/internal/proto/captcha"
	"teddy-backend/internal/repositories"
	"teddy-backend/internal/server/captcha"
	"teddy-backend/pkg/config"
	"teddy-backend/pkg/config/source/file"
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
	log.Infof("All config is %v", confType)
	if err != nil {
		log.Fatal(err)
	}
	// New Mongodb client
	mongodbClient, err := mongo.Connect(context.Background(), confType.Databases["mongodb"])
	if err != nil {
		log.Fatal(err)
	}
	// New Repository
	kvRepo, err := repositories.NewKeyValuePairRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}

	// New Handler
	captchaSrv, err := captcha.NewCaptchaServer(kvRepo)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	captchaProto.RegisterCaptchaServer(grpcServer, captchaSrv)

	healthSrv := grpcHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// Run service
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
