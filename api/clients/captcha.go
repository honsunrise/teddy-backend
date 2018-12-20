package clients

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto/captcha"
	"google.golang.org/grpc"
	"sync"
)

var captchaKey = "__teddy_captcha_client_key__"

// FromContext retrieves the client from the Context
func CaptchaFromContext(ctx *gin.Context) (captcha.CaptchaClient, bool) {
	c, ok := ctx.Value(captchaKey).(captcha.CaptchaClient)
	if c == nil {
		return nil, false
	}
	return c, ok
}

// Client returns a wrapper for the UaaClient
func CaptchaNew(f AddressFunc) gin.HandlerFunc {
	var client captcha.CaptchaClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			addr, err := f()
			if err != nil {
				log.Errorf("Get captcha address error %v", err)
				ctx.Next()
				return
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to captcha server error %v", err)
				ctx.Next()
				return
			}
			client = captcha.NewCaptchaClient(conn)
		}
		ctx.Set(captchaKey, client)
		ctx.Next()
	}
}
