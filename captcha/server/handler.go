package server

import (
	"bytes"
	"context"
	"github.com/dchest/captcha"
	"github.com/zhsyourai/teddy-backend/captcha/proto"
	"github.com/zhsyourai/teddy-backend/captcha/repositories"
	"github.com/zhsyourai/teddy-backend/common/models"
	"math/rand"
	"time"
)

const (
	Expiration = 10 * time.Minute
)

func NewCaptchaServer(repo repositories.KeyValuePairRepository) (proto.CaptchaServer, error) {
	instance := &captchaHandler{
		repo: repo,
	}
	captcha.SetCustomStore(captcha.NewMemoryStore(captcha.CollectNum, Expiration))
	return instance, nil
}

type captchaHandler struct {
	repo repositories.KeyValuePairRepository
}

func (h *captchaHandler) cleanTack() {
	for {
		<-time.After(time.Second)
		h.repo.DeleteKeyValuePairLT(time.Now())
	}
}

func (h *captchaHandler) GetCaptchaId(ctx context.Context, req *proto.GetCaptchaIdReq) (*proto.GetCaptchaIdResp, error) {
	var resp proto.GetCaptchaIdResp
	if err := validateGetCaptchaIdReq(req); err != nil {
		return nil, err
	}
	id := captcha.NewLen(int(req.Len))
	resp.Id = id
	return &resp, nil
}

func (h *captchaHandler) GetImageData(ctx context.Context, req *proto.GetImageDataReq) (*proto.GetImageDataResp, error) {
	var resp proto.GetImageDataResp
	if err := validateGetImageDataReq(req); err != nil {
		return nil, err
	}

	imgBuf := &bytes.Buffer{}

	if err := captcha.WriteImage(imgBuf, req.Id, int(req.Width), int(req.Height)); err != nil {
		return nil, err
	}
	resp.Image = imgBuf.Bytes()

	return &resp, nil
}

func (h *captchaHandler) GetVoiceData(ctx context.Context, req *proto.GetVoiceDataReq) (*proto.GetVoiceDataResp, error) {
	var resp proto.GetVoiceDataResp
	if err := validateGetVoiceDataReq(req); err != nil {
		return nil, err
	}

	voiceBuf := &bytes.Buffer{}

	if err := captcha.WriteAudio(voiceBuf, req.Id, req.Lang); err != nil {
		return nil, err
	}
	resp.VoiceWav = voiceBuf.Bytes()

	return &resp, nil
}

func (h *captchaHandler) GetRandomById(ctx context.Context, req *proto.GetRandomReq) (*proto.GetRandomResp, error) {
	var resp proto.GetRandomResp
	if err := validateGetRandomReq(req); err != nil {
		return nil, err
	}

	s := ""
	for i := 0; i < int(req.Len); i++ {
		s += (string)(rand.Intn(10) + 48)
	}

	err := h.repo.InsertKeyValuePair(&models.KeyValuePair{
		Key:        req.Id,
		Value:      s,
		ExpireTime: time.Now().Add(Expiration),
	})
	if err != nil {
		return nil, err
	}
	resp.Code = s

	return &resp, nil
}

func (h *captchaHandler) Verify(ctx context.Context, req *proto.VerifyReq) (*proto.VerifyResp, error) {
	var resp proto.VerifyResp
	if err := validateVerifyReq(req); err != nil {
		return nil, err
	}
	resp.Correct = false

	if req.Type == proto.CaptchaType_RANDOM_BY_ID {
		now := time.Now()
		_, err := h.repo.FindKeyValuePairByKeyAndValueAndExpire(req.Id, req.Code, now)
		if err != nil {
			return nil, err
		}
		resp.Correct = true
	} else if req.Type == proto.CaptchaType_IMAGE || req.Type == proto.CaptchaType_VOICE {
		if captcha.VerifyString(req.Id, req.Code) {
			resp.Correct = true
		}
	}

	return &resp, nil
}
