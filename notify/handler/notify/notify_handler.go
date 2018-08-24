package notify

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/zhsyourai/teddy-backend/notify/proto"
	"github.com/zhsyourai/teddy-backend/uaa/repositories"
	"gopkg.in/gomail.v2"
	"sync"
	"time"
)

var client *mongo.Client
var instance *accountService
var once sync.Once

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), "", clientopt.BundleClient())
	if err != nil {
		panic(err)
	}
}

func GetInstance() proto.NotifyHandler {
	once.Do(func() {
		instance = &accountService{
			repo:    repositories.NewAccountRepository(client),
			mailCh:  make(chan *gomail.Message),
			mailErr: make(chan error),
		}
		instance.startMailSender()
	})
	return instance
}

type accountService struct {
	repo    repositories.AccountRepository
	mailCh  chan *gomail.Message
	mailErr chan error
}

func (h *accountService) startMailSender() {
	go func() {
		d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")

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

func (h *accountService) SendEmail(ctx context.Context, req *proto.SendEmailReq, resp *empty.Empty) error {
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

func (h *accountService) SendInBox(ctx context.Context, req *proto.SendInBoxReq, resp *empty.Empty) error {
	panic("implement me")
}

func (h *accountService) SendNotify(ctx context.Context, req *proto.SendNotifyReq, resp *empty.Empty) error {
	panic("implement me")
}

func (h *accountService) SendSMS(ctx context.Context, req *proto.SendSMSReq, resp *empty.Empty) error {
	panic("implement me")
}

func (h *accountService) GetInBox(ctx context.Context, req *proto.GetInBoxReq, resp proto.Notify_GetInBoxStream) error {
	panic("implement me")
}

func (h *accountService) GetNotify(ctx context.Context, req *proto.GetNotifyReq, resp proto.Notify_GetNotifyStream) error {

}
