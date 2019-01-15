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
	"teddy-backend/internal/components"
	"teddy-backend/internal/mongo-grpcadapter"
	uaaProto "teddy-backend/internal/proto/uaa"
	"teddy-backend/internal/repositories"
	"teddy-backend/internal/server/uaa"
	"teddy-backend/pkg/config"
	"teddy-backend/pkg/config/source/file"
	"teddy-backend/pkg/grpcadapter"
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
	accountSrv, err := uaa.NewAccountServer(accountRepo, uidGenerator)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	uaaProto.RegisterUAAServer(grpcServer, accountSrv)

	policyAdapterServer := mongogrpcadapter.NewServer(mongodbClient, "teddy", "casbin_rule")
	grpcadapter.RegisterPolicyAdapterServer(grpcServer, policyAdapterServer)

	healthSrv := grpcHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	// Run service
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
