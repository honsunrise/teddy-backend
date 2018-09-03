package handler

import (
	"context"
	api "github.com/micro/go-api/proto"
)

type Content struct{}

// Content.Register is called by the API as /notify/inbox with post body
func (e *Content) Inbox(ctx context.Context, req *api.Request, rsp *api.Response) error {
	panic("implement me")
}

// Content.Register is called by the API as /notify/inbox with post body
func (e *Content) Notify(context.Context, *api.Request, *api.Response) error {
	panic("implement me")
}
