package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/api/uaa/handler"
	"github.com/zhsyourai/teddy-backend/common/config"
	"net/http"
	"time"
)

func main() {
	flag.Parse()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	uaa, err := handler.NewUaaHandler()
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(client.MessageNew())
	router.Use(client.UaaNew())
	router.Use(client.CaptchaNew())
	uaa.Handler(router)

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
