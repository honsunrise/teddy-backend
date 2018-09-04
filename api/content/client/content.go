package client

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"

	"github.com/zhsyourai/teddy-backend/content/proto"
)

var contentKey = "__teddy_content_client_key__"

// FromContext retrieves the client from the Context
func ContentFromContext(ctx *gin.Context) (proto.ContentService, bool) {
	c, ok := ctx.Value(contentKey).(proto.ContentService)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func ContentNew() gin.HandlerFunc {
	c := proto.NewContentService("com.teddy.srv.content", client.DefaultClient)
	return func(ctx *gin.Context) {
		ctx.Set(contentKey, c)
	}
}
