package client

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/content/proto"
	"google.golang.org/grpc"
)

var contentKey = "__teddy_content_client_key__"

// FromContext retrieves the client from the Context
func ContentFromContext(ctx *gin.Context) (proto.ContentClient, bool) {
	c, ok := ctx.Value(contentKey).(proto.ContentClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func ContentNew() gin.HandlerFunc {
	conn, err := grpc.Dial("")
	if err != nil {
		log.Errorf("Dial to captcha server error %v", err)
		return nil
	}
	client := proto.NewContentClient(conn)

	return func(ctx *gin.Context) {
		ctx.Set(contentKey, client)
		ctx.Next()
	}
}
