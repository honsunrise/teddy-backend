package client

import (
	"context"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
)

type uaaKey struct{}

// FromContext retrieves the client from the Context
func UaaFromContext(ctx context.Context) (proto.UAAService, bool) {
	c, ok := ctx.Value(uaaKey{}).(proto.UAAService)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func UaaWrapper(service micro.Service) server.HandlerWrapper {
	client := proto.NewUAAService("com.teddy.srv.uaa", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, uaaKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
