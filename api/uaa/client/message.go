package client

import (
	"context"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/zhsyourai/teddy-backend/message/proto"
)

type messageKey struct{}

// FromContext retrieves the client from the Context
func MessageFromContext(ctx context.Context) (proto.MessageService, bool) {
	c, ok := ctx.Value(messageKey{}).(proto.MessageService)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func MessageWrapper(service micro.Service) server.HandlerWrapper {
	client := proto.NewMessageService("com.teddy.srv.message", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, messageKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
