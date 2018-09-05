package message

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/xid"
	"github.com/zhsyourai/teddy-backend/common/config"
	"github.com/zhsyourai/teddy-backend/common/models"
	"github.com/zhsyourai/teddy-backend/message/converter"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"github.com/zhsyourai/teddy-backend/message/repositories"
	"gopkg.in/gomail.v2"
	"sync"
	"time"
)

func NewMessageHandler(repo repositories.InBoxRepository) (proto.MessageHandler, error) {
	instance := &notifyHandler{
		repo:    repo,
		mailCh:  make(chan *gomail.Message),
		mailErr: make(chan error),
	}
	instance.startMailSender()
	return instance, nil
}

type notifyHandler struct {
	repo    repositories.InBoxRepository
	mailCh  chan *gomail.Message
	mailErr chan error

	inBoxChMap  sync.Map
	notifyChMap sync.Map
}

func (h *notifyHandler) startMailSender() {
	// Load config
	mailConf := config.GetConfig().Mail

	go func() {
		d := gomail.NewPlainDialer(mailConf.Host, mailConf.Port, mailConf.Username, mailConf.Password)

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

func (h *notifyHandler) SendEmail(ctx context.Context, req *proto.SendEmailReq, resp *empty.Empty) error {
	if err := validateSendEmailReq(req); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", req.Topic)
	m.SetBody("text/html", req.Content)

	select {
	case err := <-h.mailErr:
		return err
	case h.mailCh <- m:
		return nil
	}
}

func (h *notifyHandler) SendSMS(ctx context.Context, req *proto.SendSMSReq, resp *empty.Empty) error {
	if err := validateSendSMSReq(req); err != nil {
		return err
	}
	return nil
}

func (h *notifyHandler) SendInBox(ctx context.Context, req *proto.SendInBoxReq, resp *empty.Empty) error {
	if err := validateSendInBoxReq(req); err != nil {
		return err
	}

	inboxItem := &models.InBoxItem{
		Unread:   true,
		ID:       xid.New().String(),
		From:     req.From,
		Type:     models.InBoxType(req.Type),
		Topic:    req.Topic,
		Content:  req.Content,
		SendTime: time.Unix(req.SendTime.Seconds, int64(req.SendTime.Nanos)),
	}
	err := h.repo.InsertInBoxItem(req.Uid, inboxItem)
	if err != nil {
		return err
	}

	if inBoxChs, ok := h.inBoxChMap.Load(req.Uid); ok {
		inBoxChs.(*sync.Map).Range(func(key, value interface{}) bool {
			key.(chan *models.InBoxItem) <- inboxItem
			return true
		})
	}
	return nil
}

func (h *notifyHandler) SendNotify(ctx context.Context, req *proto.SendNotifyReq, resp *empty.Empty) error {
	if err := validateSendNotifyReq(req); err != nil {
		return err
	}

	if inBoxChs, ok := h.notifyChMap.Load(req.Uid); ok {
		inBoxChs.(*sync.Map).Range(func(key, value interface{}) bool {
			key.(chan *models.NotifyItem) <- &models.NotifyItem{
				Uid:    req.Uid,
				Topic:  req.Topic,
				Detail: req.Detail,
			}
			return true
		})
	}
	return nil
}

func (h *notifyHandler) GetInBox(ctx context.Context, req *proto.GetInBoxReq, resp proto.Message_GetInBoxStream) error {
	if err := validateGetInBoxReq(req); err != nil {
		return err
	}
	items, err := h.repo.FindInBoxItems(req.Uid, models.InBoxType(req.Type), req.Page, req.Size, nil)
	if err != nil {
		return err
	}
	for _, item := range items {
		var pbItem proto.InBoxItem
		converter.CopyFromInBoxItemToPBInBoxItem(&item, &pbItem)
		resp.Send(&pbItem)
	}
	tmp, _ := h.inBoxChMap.LoadOrStore(req.Uid, &sync.Map{})
	inBoxChs := tmp.(*sync.Map)
	go func() {
		ch := make(chan *models.InBoxItem)
		inBoxChs.Store(ch, nil)
		var pbItem proto.InBoxItem
		for {
			item := <-ch
			converter.CopyFromInBoxItemToPBInBoxItem(item, &pbItem)
			if err := resp.Send(&pbItem); err != nil {
				resp.Close()
				close(ch)
				inBoxChs.Delete(ch)
				return
			}
		}
	}()

	return nil
}

func (h *notifyHandler) GetNotify(ctx context.Context, req *proto.GetNotifyReq, resp proto.Message_GetNotifyStream) error {
	if err := validateGetNotifyReq(req); err != nil {
		return err
	}

	tmp, _ := h.notifyChMap.LoadOrStore(req.Uid, &sync.Map{})
	inBoxChs := tmp.(*sync.Map)
	go func() {
		ch := make(chan *models.NotifyItem)
		inBoxChs.Store(ch, nil)
		var pbItem proto.NotifyItem
		for {
			item := <-ch
			converter.CopyFromNotifyItemToPBNotifyItem(item, &pbItem)
			if err := resp.Send(&pbItem); err != nil {
				resp.Close()
				close(ch)
				inBoxChs.Delete(ch)
				return
			}
		}
	}()

	return nil
}
