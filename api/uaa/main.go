package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/casbin/casbin"
	"github.com/casbin/mongodb-adapter"
	"github.com/micro/go-log"
	"github.com/zhsyourai/teddy-backend/api/uaa/handler"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/jwt"
	"github.com/zhsyourai/teddy-backend/common/utils"
	"time"

	"github.com/micro/go-micro"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/api/uaa/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.uaa"),
		micro.Version("latest"),
	)

	// Load config
	conf := config.GetConfig()
	mongodbUri := utils.BuildMongodbURI(conf.Databases["mongodb"])
	model := casbin.NewModel(conf.Casbin)

	enforcer := casbin.NewEnforcer(model, mongodbadapter.NewAdapter(mongodbUri))
	enforcer.LoadPolicy()

	// Initialise service
	service.Init(
		// create wrap for the Message srv client
		micro.WrapHandler(client.UaaWrapper(service)),
	)

	// New jwt generator and extractor
	const SigningAlgorithm = "RS512"
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal(err)
	}

	jwtGen, err := jwt.NewJwtGenerator(jwt.GeneratorConfig{
		Issuer:           "com.teddy.uaa",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return key
		},
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	jwtExt, err := jwt.NewJwtExtractor(jwt.ExtractorConfig{
		Realm:            "com.teddy",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return &key.PublicKey
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	uaa, err := handler.NewUaaHandler(enforcer, jwtGen, jwtExt)
	if err != nil {
		log.Fatal(err)
	}

	// Register Handler
	proto.RegisterUaaHandler(service.Server(), uaa)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
