package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/api/gin_jwt"
	"github.com/zhsyourai/teddy-backend/api/uaa/handler"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/config/source/file"
	"github.com/zhsyourai/teddy-backend/common/types"
	"net/http"
	"time"
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

	// New jwt generator and extractor
	const SigningAlgorithm = "RS512"
	// Load Jwt PublicKey
	block, _ := pem.Decode([]byte(confType.JWTPkcs8))
	if block == nil {
		log.Fatal("Jwt private key decode error")
	}
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	key := parseResult.(*rsa.PrivateKey)

	jwtGenerator, err := gin_jwt.NewGinJwtGenerator(gin_jwt.GeneratorConfig{
		Issuer:           "com.teddy.uaa",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return key
		},
	})
	uaa, err := handler.NewUaaHandler(jwtGenerator)
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(clients.MessageNew())
	router.Use(clients.UaaNew())
	router.Use(clients.CaptchaNew())
	uaa.Handler(router.Group("/uaa"))

	srv := http.Server{
		Addr:           fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
