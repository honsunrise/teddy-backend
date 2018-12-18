package clients

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"google.golang.org/grpc"
	"sync"
)

var uaaKey = "__teddy_uaa_client_key__"

// FromContext retrieves the client from the Context
func UaaFromContext(ctx *gin.Context) (proto.UAAClient, bool) {
	c, ok := ctx.Value(uaaKey).(proto.UAAClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func UaaNew(f AddressFunc) gin.HandlerFunc {
	var client proto.UAAClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()

			addr, err := f()
			if err != nil {
				log.Errorf("Get uaa address error %v", err)
				ctx.Next()
				return
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to uaa server error %v", err)
				ctx.Next()
				return
			}
			client = proto.NewUAAClient(conn)
		}
		ctx.Set(uaaKey, client)
		ctx.Next()
	}
}
