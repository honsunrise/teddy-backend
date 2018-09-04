package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
	"github.com/zhsyourai/teddy-backend/api/message/client"
	"github.com/zhsyourai/teddy-backend/api/message/handler"
	"github.com/zhsyourai/teddy-backend/common/config"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/micro/go-web"
)

func main() {
	flag.Parse()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	// Load config
	conf := config.GetConfig()

	// Load Jwt PublicKey
	block, _ := pem.Decode([]byte(conf.JWTPkcs8))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)

	// Create service
	service := web.NewService(
		web.Name("go.micro.api.message"),
		web.Version("latest"),
	)

	// Initialise service
	service.Init()

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:            "com.teddy.uaa",
		SigningAlgorithm: "RS512",
		KeyFunc: func() interface{} {
			return &key.PublicKey
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful handler (using Gin)
	router := gin.Default()
	router.Use(client.MessageNew())
	handler.NewMessageHandler(jwtMiddleware).Handler(router)

	// Register Handler
	service.Handle("/", router)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
