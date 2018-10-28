package client

import (
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
)

var uaaKey = "__teddy_uaa_client_key__"

// FromContext retrieves the client from the Context
func UaaFromContext(ctx *gin.Context) (proto.UAAClient, bool) {
	c, ok := ctx.Value(uaaKey).(proto.UAAClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func UaaNew() gin.HandlerFunc {
	c := proto.NewUAAClient()
	return func(ctx *gin.Context) {
		ctx.Set(uaaKey, c)
	}
}
