package clients

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"sync"
	"teddy-backend/internal/handler/errors"
	"teddy-backend/internal/proto/uaa"
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
				errors.AbortWithErrorJSON(ctx, errors.ErrGRPCDial)
				return
			}
			client = uaa.NewUAAClient(conn)
		}
		ctx.Set(uaaKey, client)
	}
}
