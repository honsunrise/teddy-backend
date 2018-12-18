package clients

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"google.golang.org/grpc"
	"sync"
)

var messageKey = "__teddy_message_client_key__"

// FromContext retrieves the client from the Context
func MessageFromContext(ctx *gin.Context) (proto.MessageClient, bool) {
	c, ok := ctx.Value(messageKey).(proto.MessageClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func MessageNew(f AddressFunc) gin.HandlerFunc {
	var client proto.MessageClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			addr, err := f()
			if err != nil {
				log.Errorf("Get message address error %v", err)
				ctx.Next()
				return
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to message server error %v", err)
				ctx.Next()
				return
			}
			client = proto.NewMessageClient(conn)
		}
		ctx.Set(messageKey, client)
		ctx.Next()
	}
}
