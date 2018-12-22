package server

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"github.com/zhsyourai/teddy-backend/content/models"
)

func copyFromTagToPBTag(tag *models.Tag, pbtag *content.Tag) error {
	if tag == nil || pbtag == nil {
		return nil
	}
	pbtag.Tag = tag.Tag
	pbtag.Usage = tag.Usage
	pbtag.Type = tag.Type
	tmp, err := ptypes.TimestampProto(tag.CreateTime)
	if err != nil {
		return err
	}
	pbtag.CreateTime = tmp

	tmp, err = ptypes.TimestampProto(tag.LastUseTime)
	if err != nil {
		return err
	}
	pbtag.LastUseTime = tmp
	return nil
}
