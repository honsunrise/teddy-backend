package client

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"google.golang.org/grpc"
)

var messageKey = "__teddy_message_client_key__"

// FromContext retrieves the client from the Context
func MessageFromContext(ctx *gin.Context) (proto.MessageClient, bool) {
	c, ok := ctx.Value(messageKey).(proto.MessageClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func MessageNew() gin.HandlerFunc {
	conn, err := grpc.Dial("srv-message")
	if err != nil {
		log.Errorf("Dial to captcha server error %v", err)
		return nil
	}
	client := proto.NewMessageClient(conn)

	return func(ctx *gin.Context) {
		ctx.Set(messageKey, client)
		ctx.Next()
	}
}
