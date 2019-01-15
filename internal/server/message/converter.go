package message

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"teddy-backend/internal/models"
	"teddy-backend/internal/proto/message"
)

func copyFromInBoxItemToPBInBoxItem(item *models.InBoxItem, pbitem *message.InBoxItem) {
	if item == nil || pbitem == nil {
		return
	}
	pbitem.Id = item.ID
	pbitem.From = item.From
	pbitem.Type = uint32(item.Type)
	pbitem.Topic = item.Topic
	pbitem.Content = item.Content
	pbitem.Unread = item.Unread
	pbitem.SendTime = &timestamp.Timestamp{
		Seconds: item.SendTime.Unix(),
		Nanos:   int32(item.SendTime.Nanosecond()),
	}
	pbitem.ReadTime = &timestamp.Timestamp{
		Seconds: item.ReadTime.Unix(),
		Nanos:   int32(item.ReadTime.Nanosecond()),
	}
}

func copyFromNotifyItemToPBNotifyItem(item *models.NotifyItem, pbitem *message.NotifyItem) {
	if item == nil || pbitem == nil {
		return
	}
	pbitem.Topic = item.Topic
	pbitem.Detail = item.Detail
}
