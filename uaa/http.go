package http

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/common/gin-jwt"
	"github.com/zhsyourai/teddy-backend/uaa/controllers"
	"net/http"
)

var (
	api             *http.Server
	webs            *http.Server
	w               errgroup.Group
	shutdownTimeout = 10 * time.Second
)

func apiServer() (*http.Server, error) {
	router := gin.Default()
	secureConf := secure.New(secure.Config{
		AllowedHosts:          []string{"example.com", "ssl.example.com"},
		SSLRedirect:           true,
		SSLHost:               "ssl.example.com",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		IsDevelopment:         false,
	})

	const SigningAlgorithm = "RS512"
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	jwtMiddleware, err := gin_jwt.NewGinJwtMiddleware(gin_jwt.MiddlewareConfig{
		Realm:            "urcf",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return &key.PublicKey
		},
	})

	if err != nil {
		return nil, err
	}

	jwtGenerator, err := gin_jwt.NewGinJwtGenerator(gin_jwt.GeneratorConfig{
		Issuer:           "urcf",
		SigningAlgorithm: SigningAlgorithm,
		KeyFunc: func() interface{} {
			return key
		},
	})

	if err != nil {
		return nil, err
	}

	router.Use(secureConf)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "DELETE"}
	corsConfig.AllowHeaders = []string{"Authorization", "Origin", "Content-Length", "Content-Type"}
	router.Use(cors.New(corsConfig))

	v1 := router.Group("/v1")
	{
		controllers.NewAccountController(jwtMiddleware, jwtGenerator).Handler(v1.Group("/uaa"))
	}

	return &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}, nil
}

func StartHTTPServer() (err error) {
	api, err = apiServer()
	if err != nil {
		log.Fatal(err)
	}

	err = api.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func main() {
	StartHTTPServer()
}
