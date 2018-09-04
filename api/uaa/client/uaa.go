package client

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"

	"github.com/zhsyourai/teddy-backend/uaa/proto"
)

var uaaKey = "__teddy_uaa_client_key__"

// FromContext retrieves the client from the Context
func UaaFromContext(ctx *gin.Context) (proto.UAAService, bool) {
	c, ok := ctx.Value(uaaKey).(proto.UAAService)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func UaaNew() gin.HandlerFunc {
	c := proto.NewUAAService("com.teddy.srv.uaa", client.DefaultClient)
	return func(ctx *gin.Context) {
		ctx.Set(uaaKey, c)
	}
}
