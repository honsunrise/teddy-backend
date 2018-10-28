package main

import (
	"flag"
	"github.com/zhsyourai/teddy-backend/api/message/client"
	"github.com/zhsyourai/teddy-backend/api/message/handler"
	"github.com/zhsyourai/teddy-backend/common/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	message, err := handler.NewMessageHandler()
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(client.MessageNew())
	message.Handler(router)

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
