package clients

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/common/proto/captcha"
	"sync"
)

var captchaKey = "__teddy_captcha_client_key__"

// FromContext retrieves the client from the Context
func CaptchaFromContext(ctx *gin.Context) captcha.CaptchaClient {
	return ctx.Value(captchaKey).(captcha.CaptchaClient)
}

// Client returns a wrapper for the UaaClient
func CaptchaNew(addr string, srv bool) gin.HandlerFunc {
	var client captcha.CaptchaClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			conn, err := getGRPCConn(addr, srv)
			if err != nil {
				ctx.Error(err)
				return
			}
			client = captcha.NewCaptchaClient(conn)
		}
		ctx.Set(messageKey, client)
		ctx.Next()
	}
}
