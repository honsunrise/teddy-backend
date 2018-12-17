package converter

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"github.com/zhsyourai/teddy-backend/content/models"
)

func CopyFromTagToPBTag(tag *models.Tag, pbtag *proto.Tag) {
	if tag == nil || pbtag == nil {
		return
	}
	pbtag.Tag = tag.Tag
	pbtag.Usage = tag.Usage
	pbtag.CreateTime = &timestamp.Timestamp{
		Seconds: tag.CreateTime.Unix(),
		Nanos:   int32(tag.CreateTime.Nanosecond()),
	}
	pbtag.LastUseTime = &timestamp.Timestamp{
		Seconds: tag.LastUseTime.Unix(),
		Nanos:   int32(tag.LastUseTime.Nanosecond()),
	}
}
