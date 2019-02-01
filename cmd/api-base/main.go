package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"teddy-backend/internal/clients"
	"teddy-backend/internal/gin_jwt"
	"teddy-backend/internal/handler/base"
	"teddy-backend/internal/handler/errors"
	"teddy-backend/pkg/config"
	"teddy-backend/pkg/config/source/file"
	"teddy-backend/pkg/grpcadapter"
	"time"
)

var (
	g errgroup.Group
)

const captchaSrvDomain = "dns:///srv-captcha:9090"
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

	adapter, err := grpcadapter.NewAdapter(uaaSrvDomain)
	if err != nil {
		log.Fatal(err)
	}

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:   "base.teddy.com",
		Issuer:  "uaa@teddy.com",
		KeyFunc: gin_jwt.RemoteFetchFunc("http:///api-uaa:8083/v1/anon/uaa/jwks.json", 24*time.Hour),
		Audience: []string{
			"base",
		},
		ErrorHandler: func(ctx *gin.Context, err error) {
			ctx.Header("WWW-Authenticate", "JWT realm=base.teddy.com")
			if err == gin_jwt.ErrForbidden {
				errors.AbortWithErrorJSON(ctx, errors.ErrForbidden)
			} else if err == gin_jwt.ErrTokenInvalid {
				errors.AbortWithErrorJSON(ctx, errors.ErrUnauthorized)
			} else if err == gin_jwt.ErrInvalidKey {
				errors.AbortWithErrorJSON(ctx, errors.ErrUnknown)
			} else {
				errors.AbortWithErrorJSON(ctx, errors.ErrUnknown)
			}
		},
	}, adapter)
	if err != nil {
		log.Fatal(err)
	}

	baseHandler, err := base.NewBaseHandler(jwtMiddleware)
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           24 * time.Hour,
	}))
	router.Use(clients.CaptchaNew(captchaSrvDomain))
	baseHandler.HandlerNormal(router.Group("/v1/anon/base").Use(jwtMiddleware.Handler()))
	baseHandler.HandlerAuth(router.Group("/v1/auth/base").Use(jwtMiddleware.Handler()))
	baseHandler.HandlerHealth(router)

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
