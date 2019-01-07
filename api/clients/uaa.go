package clients

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/common/proto/uaa"
	"google.golang.org/grpc"
	"sync"
)

var uaaKey = "__teddy_uaa_client_key__"

// FromContext retrieves the client from the Context
func UaaFromContext(ctx *gin.Context) uaa.UAAClient {
	return ctx.Value(uaaKey).(uaa.UAAClient)
}

// Client returns a wrapper for the UaaClient
func UaaNew(addr string) gin.HandlerFunc {
	var client uaa.UAAClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				ctx.Error(err)
				return
			}
			client = uaa.NewUAAClient(conn)
		}
		ctx.Set(uaaKey, client)
		ctx.Next()
	}
}
