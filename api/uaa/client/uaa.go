package client

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
	"google.golang.org/grpc"
)

var uaaKey = "__teddy_uaa_client_key__"

// FromContext retrieves the client from the Context
func UaaFromContext(ctx *gin.Context) (proto.UAAClient, bool) {
	c, ok := ctx.Value(uaaKey).(proto.UAAClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func UaaNew() gin.HandlerFunc {
	conn, err := grpc.Dial("")
	if err != nil {
		log.Errorf("Dial to captcha server error %v", err)
		return nil
	}
	client := proto.NewUAAClient(conn)

	return func(ctx *gin.Context) {
		ctx.Set(uaaKey, client)
		ctx.Next()
	}
}
