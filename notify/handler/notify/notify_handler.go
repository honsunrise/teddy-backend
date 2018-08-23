package notify

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/zhsyourai/teddy-backend/notify/proto"
	"github.com/zhsyourai/teddy-backend/uaa/repositories"
	"sync"
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
			repo: repositories.NewAccountRepository(client),
		}
	})
	return instance
}

type accountService struct {
	repo repositories.AccountRepository
}

func (h *accountService) SendEmail(context.Context, *proto.SendEmailRequest, *empty.Empty) error {
	panic("implement me")
}

func (h *accountService) SendInBox(context.Context, *proto.SendInBoxRequest, *empty.Empty) error {
	panic("implement me")
}

func (h *accountService) SendNotify(context.Context, *proto.SendNotifyRequest, *empty.Empty) error {
	panic("implement me")
}
