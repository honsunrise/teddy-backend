package client

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"

	"github.com/zhsyourai/teddy-backend/message/proto"
)

var messageKey = "__teddy_message_client_key__"

// FromContext retrieves the client from the Context
func MessageFromContext(ctx *gin.Context) (proto.MessageService, bool) {
	c, ok := ctx.Value(messageKey).(proto.MessageService)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func MessageNew() gin.HandlerFunc {
	c := proto.NewMessageService("com.teddy.srv.notify", client.DefaultClient)
	return func(ctx *gin.Context) {
		ctx.Set(messageKey, c)
	}
}
