package content

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zhsyourai/teddy-backend/content/proto"
	"github.com/zhsyourai/teddy-backend/content/repositories"
	"gopkg.in/gomail.v2"
)

func NewContentHandler(repo repositories.ContentRepository) (proto.ContentHandler, error) {
	instance := &contentHandler{
		repo:    repo,
		mailCh:  make(chan *gomail.Message),
		mailErr: make(chan error),
	}
	return instance, nil
}

type contentHandler struct {
	repo    repositories.ContentRepository
	mailCh  chan *gomail.Message
	mailErr chan error
}

func (h *contentHandler) SendEmail(context.Context, *proto.SendEmailReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) SendInBox(context.Context, *proto.SendInBoxReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) SendNotify(context.Context, *proto.SendNotifyReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) SendSMS(context.Context, *proto.SendSMSReq, *empty.Empty) error {
	panic("implement me")
}

func (h *contentHandler) GetInBox(context.Context, *proto.GetInBoxReq, proto.Content_GetInBoxStream) error {
	panic("implement me")
}

func (h *contentHandler) GetNotify(context.Context, *proto.GetNotifyReq, proto.Content_GetNotifyStream) error {
	panic("implement me")
}
