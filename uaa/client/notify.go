package client

import (
	"context"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	example "github.com/zhsyourai/teddy-backend/notify/proto"
)

type exampleKey struct{}

// FromContext retrieves the client from the Context
func NotifyFromContext(ctx context.Context) (example.NotifyService, bool) {
	c, ok := ctx.Value(exampleKey{}).(example.NotifyService)
	return c, ok
}

// Client returns a wrapper for the NotifyClient
func NotifyWrapper(service micro.Service) server.HandlerWrapper {
	client := example.NewNotifyService("com.teddy.srv.notify", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, exampleKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
