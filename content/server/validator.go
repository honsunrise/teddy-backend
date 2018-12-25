package server

import "github.com/zhsyourai/teddy-backend/common/proto/content"

func validatePublishInfoReq(req *content.PublishInfoReq) error {
	if req.Uid == "" {

	} else if req.Author == "" {

	} else if req.Title == "" {

	} else if len(req.Tags) < 1 {

	} else if req.Summary == "" {

	} else if len(req.CoverResources) < 1 {

	} else if req.Country == "" {

	}
	return nil
}

func validateGetTagsReq(req *content.GetTagReq) error {
	return nil
}

func validateEditInfoReq(req *content.EditInfoReq) error {
	return nil
}

func validateGetInfoReq(req *content.GetInfoReq) error {
	return nil
}

func validateGetInfosReq(req *content.GetInfosReq) error {
	return nil
}

func validateInfoIDAndUIDReq(req *content.InfoIDAndUIDReq) error {
	return nil
}

func validateUIDPageReq(req *content.UIDPageReq) error {
	return nil
}

func validateInfoIDPageReq(req *content.InfoIDPageReq) error {
	return nil
}

func validateGetSegmentsReq(req *content.GetSegmentsReq) error {
	return nil
}

func validatePublishSegmentReq(req *content.PublishSegmentReq) error {
	return nil
}

func validateEditSegmentReq(req *content.EditSegmentReq) error {
	return nil
}

func validateInfoIDAndUIDAndSegIDReq(req *content.InfoIDAndUIDAndSegIDReq) error {
	return nil
}
