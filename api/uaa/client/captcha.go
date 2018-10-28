package client

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/captcha/proto"
)

var captchaKey = "__teddy_captcha_client_key__"

// FromContext retrieves the client from the Context
func CaptchaKeyFromContext(ctx *gin.Context) (proto.CaptchaClient, bool) {
	c, ok := ctx.Value(captchaKey).(proto.CaptchaClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func CaptchaNew() gin.HandlerFunc {
	c := proto.NewCaptchaClient()
	return func(ctx *gin.Context) {
		ctx.Set(captchaKey, c)
	}
}
