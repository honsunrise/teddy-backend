package server

import (
	"github.com/zhsyourai/teddy-backend/common/proto/captcha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateGetCaptchaIdReq(req *captcha.GetCaptchaIdReq) error {
	if req.Len < 4 {
		return status.Error(codes.InvalidArgument, "captcha len must gte 4")
	}
	return nil
}

func validateGetImageDataReq(req *captcha.GetImageDataReq) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "image captcha id must not be empty")
	} else if req.Height < 66 || req.Height > 300 {
		return status.Error(codes.InvalidArgument, "image captcha height must gte 100 and lte 300")
	} else if req.Width < 200 || req.Width > 400 {
		return status.Error(codes.InvalidArgument, "image captcha height must gte 100 and lte 300")
	}
	return nil
}

func validateGetRandomReq(req *captcha.GetRandomReq) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "captcha id must not be empty")
	}
	return nil
}

func validateGetVoiceDataReq(req *captcha.GetVoiceDataReq) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "voice captcha id must not be empty")
	} else if req.Lang == "" {
		return status.Error(codes.InvalidArgument, "voice captcha lang must not be empty")
	}
	return nil
}

func validateVerifyReq(req *captcha.VerifyReq) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "captcha id must not be empty")
	} else if req.Code == "" {
		return status.Error(codes.InvalidArgument, "captcha code must not be empty")
	}
	return nil
}
