package clients

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var messageKey = "__teddy_message_client_key__"

const messageSrvDomain = "srv-message"

// FromContext retrieves the client from the Context
func MessageFromContext(ctx *gin.Context) (proto.MessageClient, bool) {
	c, ok := ctx.Value(messageKey).(proto.MessageClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func MessageNew() gin.HandlerFunc {
	var client proto.MessageClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			_, addrs, err := net.LookupSRV("http", "tcp", messageSrvDomain)
			if err != nil {
				log.Errorf("Dial to captcha server error %v", err)
				ctx.Next()
				return
			}
			for _, addr := range addrs {
				log.Infof("%s SRV is %v", messageSrvDomain, addr)
			}
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", messageSrvDomain, addrs[0].Port), grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to captcha server error %v", err)
				ctx.Next()
				return
			}
			client = proto.NewMessageClient(conn)
		}
		ctx.Set(messageKey, client)
		ctx.Next()
	}
}
