package client

import (
	"github.com/gin-gonic/gin"

	"github.com/zhsyourai/teddy-backend/message/proto"
)

var messageKey = "__teddy_message_client_key__"

// FromContext retrieves the client from the Context
func MessageFromContext(ctx *gin.Context) (proto.MessageClient, bool) {
	c, ok := ctx.Value(messageKey).(proto.MessageClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func MessageNew() gin.HandlerFunc {
	c := proto.NewMessageClient()
	return func(ctx *gin.Context) {
		ctx.Set(messageKey, c)
	}
}
