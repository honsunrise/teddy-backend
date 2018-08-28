package message

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"github.com/zhsyourai/teddy-backend/message/repositories"
	"gopkg.in/gomail.v2"
	"time"
)

func NewMessageHandler(repo repositories.InBoxRepository) (proto.MessageHandler, error) {
	instance := &notifyService{
		repo:    repo,
		mailCh:  make(chan *gomail.Message),
		mailErr: make(chan error),
	}
	instance.startMailSender()
	return instance, nil
}

type notifyService struct {
	repo    repositories.InBoxRepository
	mailCh  chan *gomail.Message
	mailErr chan error
}

func (h *notifyService) startMailSender() {
	go func() {
		d := gomail.NewPlainDialer("smtp.example.com", 587, "user", "123456")

		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-h.mailCh:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						h.mailErr <- err
						open = false
					} else {
						open = true
					}
				} else {
					if err := gomail.Send(s, m); err != nil {
						h.mailErr <- err
					}
				}
			case <-time.After(30 * time.Second):
				if open {
					s.Close()
					open = false
				}
			}
		}
	}()
}

func (h *notifyService) SendEmail(ctx context.Context, req *proto.SendEmailReq, resp *empty.Empty) error {
	if err := validateSendEmailReq(req); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", req.Topic)
	m.SetBody("text/html", req.Cotent)

	select {
	case err := <-h.mailErr:
		return err
	case h.mailCh <- m:
		return nil
	}
}

func (h *notifyService) SendInBox(ctx context.Context, req *proto.SendInBoxReq, resp *empty.Empty) error {
	if err := validateSendInBoxReq(req); err != nil {
		return err
	}
	return nil
}

func (h *notifyService) SendNotify(ctx context.Context, req *proto.SendNotifyReq, resp *empty.Empty) error {
	if err := validateSendNotifyReq(req); err != nil {
		return err
	}
	return nil
}

func (h *notifyService) SendSMS(ctx context.Context, req *proto.SendSMSReq, resp *empty.Empty) error {
	if err := validateSendSMSReq(req); err != nil {
		return err
	}
	return nil
}

func (h *notifyService) GetInBox(ctx context.Context, req *proto.GetInBoxReq, resp proto.Message_GetInBoxStream) error {
	if err := validateGetInBoxReq(req); err != nil {
		return err
	}
	return nil
}

func (h *notifyService) GetNotify(ctx context.Context, req *proto.GetNotifyReq, resp proto.Message_GetNotifyStream) error {
	if err := validateGetNotifyReq(req); err != nil {
		return err
	}
	return nil
}
