package main

import (
	"errors"
	"fmt"
	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/config/source/file"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"github.com/zhsyourai/teddy-backend/content/server"
	"google.golang.org/grpc"
	grpcHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"

	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/config"
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

	minioConfig := confType.ObjectStore["minio"]
	if minioConfig == nil {
		log.Fatal(errors.New("missing minio config"))
	}

	// Initialize minio client object.
	minioClient, err := minio.New(minioConfig.Endpoint, minioConfig.AccessKey, minioConfig.SecretKey, false)
	if err != nil {
		log.Fatal(err)
	}

	// New Handler
	accountHandler, err := server.NewContentServer(mongodbClient, minioClient)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	content.RegisterContentServer(grpcServer, accountHandler)

	healthSrv := grpcHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// Run service
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
