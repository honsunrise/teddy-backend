package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"github.com/zhsyourai/teddy-backend/api/content/client"
	"github.com/zhsyourai/teddy-backend/api/content/handler"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
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
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	key := parseResult.(*rsa.PrivateKey)

	// Create service
	service := web.NewService(
		web.Name("go.micro.api.content"),
		web.Version("latest"),
	)

	// Initialise service
	service.Init()

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:            "com.teddy",
		SigningAlgorithm: "RS512",
		KeyFunc: func() interface{} {
			return &key.PublicKey
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	content, err := handler.NewContentHandler(jwtMiddleware)
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful handler (using Gin)
	router := gin.Default()
	router.Use(client.ContentNew())
	content.Handler(router)

	// Register Handler
	service.Handle("/", router)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
