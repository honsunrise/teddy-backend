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

func copyFromSegmentToPBSegment(tag *models.Segment, pbtag *content.Segment) error {
	if tag == nil || pbtag == nil {
		return nil
	}
	pbtag.Id = tag.ID.Hex()
	pbtag.InfoID = tag.InfoID.Hex()
	pbtag.No = tag.No
	pbtag.Title = tag.Title
	pbtag.Labels = tag.Labels
	return nil
}

func copyFromValueToPBValue(value *models.Value, pbValue *content.Value) error {
	if value == nil || pbValue == nil {
		return nil
	}
	pbValue.Id = value.ID
	tmp, err := ptypes.TimestampProto(value.Time)
	if err != nil {
		return err
	}
	pbValue.Time = tmp
	pbValue.Value = value.Value
	return nil
}
