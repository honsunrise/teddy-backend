package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/message/converter"
	"github.com/zhsyourai/teddy-backend/message/models"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"github.com/zhsyourai/teddy-backend/message/repositories"
	"gopkg.in/gomail.v2"
	"sync"
	"time"
)

func NewMessageServer(repo repositories.InBoxRepository, host string, port int, username string, password string) (proto.MessageServer, error) {
	instance := &notifyHandler{
		repo:    repo,
		mailCh:  make(chan *gomail.Message),
		mailErr: make(chan error),
	}
	instance.startMailSender()
	return instance, nil
}

type notifyHandler struct {
	repo     repositories.InBoxRepository
	mailCh   chan *gomail.Message
	mailErr  chan error
	host     string
	port     int
	username string
	password string

	notifyChMap sync.Map
}

func (h *notifyHandler) startMailSender() {
	// Load config

	go func() {
		d := gomail.NewPlainDialer(h.host, h.port, h.username, h.password)

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

func (h *notifyHandler) SendEmail(ctx context.Context, req *proto.SendEmailReq) (*empty.Empty, error) {
	log.Infof("Send Email to %v", req)

	var resp empty.Empty

	if err := validateSendEmailReq(req); err != nil {
		return nil, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", req.Topic)
	m.SetBody("text/html", req.Content)

	select {
	case err := <-h.mailErr:
		return nil, err
	case h.mailCh <- m:
		return &resp, nil
	}
}

func (h *notifyHandler) SendSMS(ctx context.Context, req *proto.SendSMSReq) (*empty.Empty, error) {
	var resp empty.Empty

	if err := validateSendSMSReq(req); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *notifyHandler) SendInBox(ctx context.Context, req *proto.SendInBoxReq) (*empty.Empty, error) {
	var resp empty.Empty

	if err := validateSendInBoxReq(req); err != nil {
		return nil, err
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
		return nil, err
	}
	return &resp, nil
}

func (h *notifyHandler) SendNotify(ctx context.Context, req *proto.SendNotifyReq) (*empty.Empty, error) {
	var resp empty.Empty

	if err := validateSendNotifyReq(req); err != nil {
		return nil, err
	}

	if inBoxCh, ok := h.notifyChMap.Load(req.Uid); ok {
		inBoxCh.(chan *models.NotifyItem) <- &models.NotifyItem{
			Uid:    req.Uid,
			Topic:  req.Topic,
			Detail: req.Detail,
		}
	}
	return &resp, nil
}

func (h *notifyHandler) GetInBox(ctx context.Context, req *proto.GetInBoxReq) (*proto.GetInboxResp, error) {
	var resp proto.GetInboxResp
	if err := validateGetInBoxReq(req); err != nil {
		return nil, err
	}
	items, err := h.repo.FindInBoxItems(req.Uid, models.InBoxType(req.Type), req.Page, req.Size, nil)
	if err != nil {
		return nil, err
	}
	resp.Items = make([]*proto.InBoxItem, len(items))
	for i, item := range items {
		var pbItem proto.InBoxItem
		converter.CopyFromInBoxItemToPBInBoxItem(&item, &pbItem)
		resp.Items[i] = &pbItem
	}

	return &resp, nil
}

func (h *notifyHandler) GetNotify(req *proto.GetNotifyReq, resp proto.Message_GetNotifyServer) error {
	if err := validateGetNotifyReq(req); err != nil {
		return err
	}

	tmp, _ := h.notifyChMap.LoadOrStore(req.Uid, make(chan *models.NotifyItem))
	inBoxCh := tmp.(chan *models.NotifyItem)
	var pbItem proto.NotifyItem
	for {
		item := <-inBoxCh
		converter.CopyFromNotifyItemToPBNotifyItem(item, &pbItem)
		if err := resp.Send(&pbItem); err != nil {
			close(inBoxCh)
			h.notifyChMap.Delete(req.Uid)
			return err
		}
	}
	return nil
}
