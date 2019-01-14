package server

import (
	"github.com/golang/protobuf/ptypes"
	"teddy-backend/common/proto/content"
	"teddy-backend/content/models"
)

func copyFromTagToPBTag(tag *models.Tag, pbTag *content.TagResp) error {
	if tag == nil || pbTag == nil {
		return nil
	}
	pbTag.Tag = tag.Tag
	pbTag.Usage = tag.Usage
	pbTag.Type = tag.Type
	tmp, err := ptypes.TimestampProto(tag.CreateTime)
	if err != nil {
		return err
	}
	pbTag.CreateTime = tmp

	tmp, err = ptypes.TimestampProto(tag.LastUseTime)
	if err != nil {
		return err
	}
	pbTag.LastUseTime = tmp
	return nil
}

func copyFromSegmentToPBSegment(segment *models.Segment, pbSegment *content.SegmentResp) error {
	if segment == nil || pbSegment == nil {
		return nil
	}
	pbSegment.Id = segment.ID.Hex()
	pbSegment.InfoID = segment.InfoID.Hex()
	pbSegment.No = segment.No
	pbSegment.Title = segment.Title
	pbSegment.Labels = segment.Labels
	return nil
}

func copyFromValueToPBValue(value *models.Value, pbValue *content.ValueResp) error {
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
