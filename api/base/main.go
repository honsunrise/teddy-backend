package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/base/handler"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/api/gin_jwt"
	"github.com/zhsyourai/teddy-backend/api/nice_error"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/config/source/file"
	"golang.org/x/sync/errgroup"
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

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:            "uaa.teddy.com",
		Issuer:           "uaa@teddy.com",
		SigningAlgorithm: "RS256",
		KeyFunc:          gin_jwt.RemoteFetchFunc("http://api-uaa:8083/v1/anon/jwks.json", 24*time.Hour),
		Audience: []string{
			"base",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	base, err := handler.NewBaseHandler(jwtMiddleware)
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(clients.CaptchaNew(captchaSrvAddrFunc))
	router.Use(nice_error.NewNiceError())
	base.HandlerNormal(router.Group("/v1/anon/base").Use(jwtMiddleware.Handler(true)))
	base.HandlerAuth(router.Group("/v1/auth/base").Use(jwtMiddleware.Handler(false)))
	base.HandlerHealth(router)

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
