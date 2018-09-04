package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/casbin/casbin"
	"github.com/casbin/mongodb-adapter"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-log"
	"github.com/micro/go-web"
	"github.com/micro/micro/web"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/api/uaa/components"
	"github.com/zhsyourai/teddy-backend/api/uaa/handler"
	"github.com/zhsyourai/teddy-backend/api/uaa/repositories"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/utils"
)

func main() {
	// Load config
	conf := config.GetConfig()
	mongodbUri := utils.BuildMongodbURI(conf.Databases["mongodb"])
	model := casbin.NewModel(conf.Casbin)

	enforcer := casbin.NewEnforcer(model, mongodbadapter.NewAdapter(mongodbUri))
	enforcer.LoadPolicy()

	// New Mongodb client
	mongodbClient, err := mongo.Connect(context.Background(), mongodbUri)
	if err != nil {
		log.Fatal(err)
	}
	// New Repository
	keyValuePairRepo, err := repositories.NewKeyValuePairRepository(mongodbClient)
	if err != nil {
		log.Fatal(err)
	}
	simpleCaptch, err := components.NewCaptchaVerifier(keyValuePairRepo)
	if err != nil {
		log.Fatal(err)
	}
	// New Service
	service := web.NewService(
		web.Name("go.micro.api.uaa"),
		web.Version("latest"),
	)
	// Initialise service
	service.Init()

	// New jwt generator and extractor
	const SigningAlgorithm = "RS512"
	// Load Jwt PublicKey
	block, _ := pem.Decode([]byte(conf.JWTPkcs8))
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	key := parseResult.(*rsa.PrivateKey)

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:            "com.teddy",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return &key.PublicKey
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	jwtGenerator, err := gin_jwt.NewGinJwtGenerator(gin_jwt.GeneratorConfig{
		Issuer:           "com.teddy.uaa",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return key
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	uaa, err := handler.NewUaaHandler(enforcer, simpleCaptch, jwtMiddleware, jwtGenerator)
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful handler (using Gin)
	router := gin.Default()
	router.Use(client.MessageNew())
	router.Use(client.UaaNew())
	uaa.Handler(router)

	// Register Handler
	service.Handle("/", router)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
