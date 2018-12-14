package clients

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var captchaKey = "__teddy_captcha_client_key__"

const captchaSrvDomain = "srv-captcha"

// FromContext retrieves the client from the Context
func CaptchaFromContext(ctx *gin.Context) (proto.CaptchaClient, bool) {
	c, ok := ctx.Value(captchaKey).(proto.CaptchaClient)
	if c == nil {
		return nil, false
	}
	return c, ok
}

// Client returns a wrapper for the UaaClient
func CaptchaNew() gin.HandlerFunc {
	var client proto.CaptchaClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			_, addrs, err := net.LookupSRV("grpc", "tcp", captchaSrvDomain)
			if err != nil {
				log.Errorf("Dial to captcha server error %v", err)
				ctx.Next()
				return
			}
			for _, addr := range addrs {
				log.Infof("%s SRV is %v", captchaSrvDomain, addr)
			}
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", captchaSrvDomain, addrs[0].Port), grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to captcha server error %v", err)
				ctx.Next()
				return
			}
			client = proto.NewCaptchaClient(conn)
		}
		ctx.Set(captchaKey, client)
		ctx.Next()
	}
}
