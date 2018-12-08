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
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	g errgroup.Group
)

func init() {
	log.SetReportCaller(true)
}

func main() {
	conf, err := config.NewConfig(file.NewSource(file.WithFormat(config.Yaml), file.WithPath("config/config.yaml")))
	if err != nil {
		log.Fatal(err)
	}
	var confType Config
	err = conf.Scan(&confType)
	if err != nil {
		log.Fatal(err)
	}

	certPEM, err := ioutil.ReadFile("secret/JWTPkcs8")
	if err != nil {
		log.Fatal(err)
	}
	// New jwt generator and extractor
	const SigningAlgorithm = "RS512"
	// Load Jwt PublicKey
	block, _ := pem.Decode(certPEM)
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
	routerNormal := gin.Default()
	routerNormal.Use(clients.MessageNew())
	routerNormal.Use(clients.UaaNew())
	routerNormal.Use(clients.CaptchaNew())
	uaa.HandlerNormal(routerNormal.Group("/uaa"))

	routerHealth := gin.Default()
	routerNormal.Use(clients.MessageNew())
	routerNormal.Use(clients.UaaNew())
	routerNormal.Use(clients.CaptchaNew())
	uaa.HandlerHealth(routerHealth)

	// For normal request
	srv1 := http.Server{
		Addr:         fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port),
		Handler:      routerNormal,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// For health check port
	srv2 := http.Server{
		Addr:         fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port+100),
		Handler:      routerHealth,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return srv1.ListenAndServe()
	})

	g.Go(func() error {
		return srv2.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
