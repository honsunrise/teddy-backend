package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"teddy-backend/api/clients"
	"teddy-backend/api/uaa/handler"
	"teddy-backend/common/config"
	"teddy-backend/common/config/source/file"
	"teddy-backend/common/gin_jwt"
	"teddy-backend/common/grpcadapter"
	"teddy-backend/common/nice_error"
	"time"
)

var (
	g errgroup.Group
)

const captchaSrvDomain = "dns:///srv-captcha:9090"
const messageSrvDomain = "dns:///srv-message:9092"
const uaaSrvDomain = "dns:///srv-uaa:9093"

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
	const SigningAlgorithm = "RS256"
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
		Issuer:           "uaa@teddy.com",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return key
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	adapter, err := grpcadapter.NewAdapter(uaaSrvDomain)
	if err != nil {
		log.Fatal(err)
	}

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:  "uaa.teddy.com",
		Issuer: "uaa@teddy.com",
		KeyFunc: func() interface{} {
			return key.Public()
		},
		Audience: []string{
			"uaa",
		},
	}, adapter)
	if err != nil {
		log.Fatal(err)
	}
	uaa, err := handler.NewUaaHandler(jwtMiddleware, jwtGenerator)
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		// Process request
		path := ctx.Request.URL.Path
		ctx.Next()

		log.Infof("PATH %s HEADERS %v", path, ctx.Request.Header)
	})
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           24 * time.Hour,
	}))
	router.Use(clients.MessageNew(messageSrvDomain))
	router.Use(clients.UaaNew(uaaSrvDomain))
	router.Use(clients.CaptchaNew(captchaSrvDomain))
	router.Use(nice_error.NewNiceError())
	uaa.HandlerNormal(router.Group("/v1/anon/uaa").Use(jwtMiddleware.Handler()))
	uaa.HandlerAuth(router.Group("/v1/auth/uaa").Use(jwtMiddleware.Handler()))
	uaa.HandlerHealth(router)

	// For normal request
	srv1 := http.Server{
		Addr:         fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// For health check port
	srv2 := http.Server{
		Addr:         fmt.Sprintf("%s:%d", confType.Server.Address, confType.Server.Port+100),
		Handler:      router,
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
