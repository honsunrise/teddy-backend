package clients

import (
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"google.golang.org/grpc"
)

var contentKey = "__teddy_content_client_key__"

// FromContext retrieves the client from the Context
func ContentFromContext(ctx *gin.Context) (content.ContentClient, bool) {
	c, ok := ctx.Value(contentKey).(content.ContentClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func ContentNew(f AddressFunc) gin.HandlerFunc {
	var client content.ContentClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			addr, err := f()
			if err != nil {
				log.Errorf("Get content address error %v", err)
				ctx.Next()
				return
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to content server error %v", err)
				ctx.Next()
				return
			}
			client = content.NewContentClient(conn)
		}
		ctx.Set(contentKey, client)
		ctx.Next()
	}
}
