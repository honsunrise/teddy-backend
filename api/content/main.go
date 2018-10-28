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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/micro/go-web"
)

func main() {
	flag.Parse()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	content, err := handler.NewContentHandler()
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(client.ContentNew())
	content.Handler(router)

	srv := http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
