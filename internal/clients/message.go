package clients

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"sync"
	"teddy-backend/internal/handler/errors"
	"teddy-backend/internal/proto/message"
)

var messageKey = "__teddy_message_client_key__"

// FromContext retrieves the client from the Context
func MessageFromContext(ctx *gin.Context) message.MessageClient {
	return ctx.Value(messageKey).(message.MessageClient)
}

// Client returns a wrapper for the UaaClient
func MessageNew(addr string) gin.HandlerFunc {
	var client message.MessageClient = nil
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
			client = message.NewMessageClient(conn)
		}
		ctx.Set(messageKey, client)
	}
}
