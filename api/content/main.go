package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"teddy-backend/api/clients"
	"teddy-backend/api/content/handler"
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
const contentSrvDomain = "dns:///srv-content:9091"
const uaaSrvDomain = "dns:///srv-uaa:9093"

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

	adapter, err := grpcadapter.NewAdapter(uaaSrvDomain)
	if err != nil {
		log.Fatal(err)
	}

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:   "content.teddy.com",
		Issuer:  "uaa@teddy.com",
		KeyFunc: gin_jwt.RemoteFetchFunc("http://10.10.10.30:8083/v1/anon/uaa/jwks.json", 24*time.Hour),
		Audience: []string{
			"content",
		},
	}, adapter)
	if err != nil {
		log.Fatal(err)
	}

	minioConfig := confType.ObjectStore["minio"]
	if minioConfig == nil {
		log.Fatal(errors.New("missing minio config"))
	}

	// Initialize minio client object.
	minioClient, err := minio.New(minioConfig.Endpoint, minioConfig.AccessKey, minioConfig.SecretKey, false)
	if err != nil {
		log.Fatal(err)
	}

	content, err := handler.NewContentHandler(jwtMiddleware, minioClient, minioConfig.Bucket)
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
	router.Use(clients.ContentNew(contentSrvDomain))
	router.Use(clients.CaptchaNew(captchaSrvDomain))
	router.Use(nice_error.NewNiceError())
	content.HandlerNormal(router.Group("/v1/anon/content").Use(jwtMiddleware.Handler()))
	content.HandlerAuth(router.Group("/v1/auth/content").Use(jwtMiddleware.Handler()))
	content.HandlerHealth(router)

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
