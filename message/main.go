package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	grpcHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"teddy-backend/common/config/source/file"

	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"teddy-backend/common/config"
	"teddy-backend/common/proto/message"
	"teddy-backend/message/repositories"
	"teddy-backend/message/server"
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
	// New Mongodb client
	mongodbClient, err := mongo.Connect(context.Background(), confType.Databases["mongodb"])
	if err != nil {
		log.Fatal(err)
	}
	// New Repository
	inboxRepo, err := repositories.NewInBoxRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}

	// New Handler
	messageSrv, err := server.NewMessageServer(inboxRepo, confType.Mail.Host, confType.Mail.Port,
		confType.Mail.Username, confType.Mail.Password)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	message.RegisterMessageServer(grpcServer, messageSrv)

	healthSrv := grpcHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// Run service
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
