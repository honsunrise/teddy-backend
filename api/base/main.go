package main

import (
	"github.com/zhsyourai/teddy-backend/api/base/client"
	"github.com/zhsyourai/teddy-backend/api/base/handler"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/config/source/file"
	"github.com/zhsyourai/teddy-backend/common/types"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.NewConfig(file.NewSource(file.WithFormat(config.Yaml), file.WithPath("config.yaml")))
	if err != nil {
		log.Fatal(err)
	}
	var confType types.Config
	err = conf.Scan(&confType)
	if err != nil {
		log.Fatal(err)
	}

	content, err := handler.NewBaseHandler()
	if err != nil {
		log.Fatal(err)
	}

	// Create RESTful server (using Gin)
	router := gin.Default()
	router.Use(client.CaptchaNew())
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
