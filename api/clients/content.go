package clients

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"google.golang.org/grpc"
	"sync"
)

var contentKey = "__teddy_content_client_key__"

// FromContext retrieves the client from the Context
func ContentFromContext(ctx *gin.Context) content.ContentClient {
	return ctx.Value(contentKey).(content.ContentClient)
}

// Client returns a wrapper for the UaaClient
func ContentNew(addr string) gin.HandlerFunc {
	var client content.ContentClient = nil
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
			client = content.NewContentClient(conn)
		}
		ctx.Set(messageKey, client)
		ctx.Next()
	}
}
