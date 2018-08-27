package client

import (
	"context"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/zhsyourai/teddy-backend/notify/proto"
)

type notifyKey struct{}

// FromContext retrieves the client from the Context
func NotifyFromContext(ctx context.Context) (proto.NotifyService, bool) {
	c, ok := ctx.Value(notifyKey{}).(proto.NotifyService)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func NotifyWrapper(service micro.Service) server.HandlerWrapper {
	client := proto.NewNotifyService("com.teddy.srv.notify", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, notifyKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
