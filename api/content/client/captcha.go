package client

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/captcha/proto"
	"google.golang.org/grpc"
)

var captchaKey = "__teddy_captcha_client_key__"

// FromContext retrieves the client from the Context
func CaptchaFromContext(ctx *gin.Context) (proto.CaptchaClient, bool) {
	c, ok := ctx.Value(captchaKey).(proto.CaptchaClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func CaptchaNew() gin.HandlerFunc {
	conn, err := grpc.Dial("")
	if err != nil {
		log.Errorf("Dial to captcha server error %v", err)
		return nil
	}
	client := proto.NewCaptchaClient(conn)

	return func(ctx *gin.Context) {
		ctx.Set(captchaKey, client)
		ctx.Next()
	}
}
