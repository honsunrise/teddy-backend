package converter

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/zhsyourai/teddy-backend/common/models"
	"github.com/zhsyourai/teddy-backend/message/proto"
	"time"
)

func CopyFromInBoxItemToPBInBoxItem(item *models.InBoxItem, pbitem *proto.InBoxItem) {
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

func CopyFromPBInBoxItemToInBoxItem(pbitem *proto.InBoxItem, item *models.InBoxItem) {
	if item == nil || pbitem == nil {
		return
	}
	item.ID = pbitem.Id
	item.From = pbitem.From
	item.Type = models.InBoxType(pbitem.Type)
	item.Topic = pbitem.Topic
	item.Content = pbitem.Content
	item.Unread = pbitem.Unread
	item.SendTime = time.Unix(pbitem.SendTime.Seconds, int64(pbitem.SendTime.Nanos))
	item.ReadTime = time.Unix(pbitem.ReadTime.Seconds, int64(pbitem.ReadTime.Nanos))
}

func CopyFromNotifyItemToPBNotifyItem(item *models.NotifyItem, pbitem *proto.NotifyItem) {
	if item == nil || pbitem == nil {
		return
	}
	pbitem.Topic = item.Topic
	pbitem.Detail = item.Detail
}
