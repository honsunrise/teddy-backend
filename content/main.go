package main

import (
	"fmt"
	"github.com/micro/go-micro"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/utils"
	"github.com/zhsyourai/teddy-backend/content/components"
	"github.com/zhsyourai/teddy-backend/content/server"
	"google.golang.org/grpc"
	"net"

	"context"
	"flag"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/content/proto"
	"github.com/zhsyourai/teddy-backend/content/repositories"
)

const PORT = 9999

func main() {
	flag.Parse()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	// Load config
	conf := config.GetConfig()
	mongodbUri := utils.BuildMongodbURI(conf.Databases["mongodb"])

	// New Mongodb client
	mongodbClient, err := mongo.Connect(context.Background(), mongodbUri)
	if err != nil {
		log.Fatal(err)
	}
	// New Repository
	accountRepo, err := repositories.NewInfoRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}
	// New Service
	service := micro.NewService(
		micro.Name("com.teddy.srv.uaa"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()
	// New Handler
	accountHandler, err := server.NewContentServer(accountRepo)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterContentServer(grpcServer, accountHandler)

	// Run service
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
