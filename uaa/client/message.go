package client

import (
	"context"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	proto "github.com/zhsyourai/teddy-backend/message/proto"
)

type protoKey struct{}

// FromContext retrieves the client from the Context
func MessageFromContext(ctx context.Context) (proto.MessageService, bool) {
	c, ok := ctx.Value(protoKey{}).(proto.MessageService)
	return c, ok
}

// Client returns a wrapper for the NotifyClient
func NotifyWrapper(service micro.Service) server.HandlerWrapper {
	client := proto.NewMessageService("com.teddy.srv.message", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, protoKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
