package clients

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"sync"
	"teddy-backend/internal/handler/errors"
	"teddy-backend/internal/proto/captcha"
)

var captchaKey = "__teddy_captcha_client_key__"

// FromContext retrieves the client from the Context
func CaptchaFromContext(ctx *gin.Context) captcha.CaptchaClient {
	return ctx.Value(captchaKey).(captcha.CaptchaClient)
}

// Client returns a wrapper for the UaaClient
func CaptchaNew(addr string) gin.HandlerFunc {
	var client captcha.CaptchaClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				errors.AbortWithErrorJSON(ctx, errors.ErrGRPCDial)
				return
			}
			client = captcha.NewCaptchaClient(conn)
		}
		ctx.Set(captchaKey, client)
	}
}
