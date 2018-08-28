package client

import (
	"context"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/zhsyourai/teddy-backend/message/proto"
)

type notifyKey struct{}

// FromContext retrieves the client from the Context
func MessageFromContext(ctx context.Context) (proto.MessageService, bool) {
	c, ok := ctx.Value(notifyKey{}).(proto.MessageService)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func MessageWrapper(service micro.Service) server.HandlerWrapper {
	client := proto.NewMessageService("com.teddy.srv.notify", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, notifyKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
