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

func (h *captchaHandler) GetCaptchaId(ctx context.Context, req *proto.GetCaptchaIdReq, rsp *proto.GetCaptchaIdResp) error {
	if err := validateGetCaptchaIdReq(req); err != nil {
		return err
	}
	id := captcha.NewLen(int(req.Len))
	rsp.Id = id
	return nil
}

func (h *captchaHandler) GetImageData(ctx context.Context, req *proto.GetImageDataReq, rsp *proto.GetImageDataResp) error {
	if err := validateGetImageDataReq(req); err != nil {
		return err
	}

	imgBuf := &bytes.Buffer{}

	if err := captcha.WriteImage(imgBuf, req.Id, int(req.Width), int(req.Height)); err != nil {
		return err
	}
	rsp.Image = imgBuf.Bytes()

	return nil
}

func (h *captchaHandler) GetVoiceData(ctx context.Context, req *proto.GetVoiceDataReq, rsp *proto.GetVoiceDataResp) error {
	if err := validateGetVoiceDataReq(req); err != nil {
		return err
	}

	voiceBuf := &bytes.Buffer{}

	if err := captcha.WriteAudio(voiceBuf, req.Id, req.Lang); err != nil {
		return err
	}
	rsp.VoiceWav = voiceBuf.Bytes()

	return nil
}

func (h *captchaHandler) GetRandomById(ctx context.Context, req *proto.GetRandomReq, rsp *proto.GetRandomResp) error {
	if err := validateGetRandomReq(req); err != nil {
		return err
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
		return err
	}
	return nil
}

func (h *captchaHandler) Verify(ctx context.Context, req *proto.VerifyReq, rsp *proto.VerifyResp) error {
	if err := validateVerifyReq(req); err != nil {
		return err
	}
	rsp.Correct = false

	if req.Type == proto.CaptchaType_RANDOM_BY_ID {
		now := time.Now()
		_, err := h.repo.FindKeyValuePairByKeyAndValueAndExpire(req.Id, req.Code, now)
		if err != nil {
			return err
		}
		rsp.Correct = true
	} else if req.Type == proto.CaptchaType_IMAGE || req.Type == proto.CaptchaType_VOICE {
		if captcha.VerifyString(req.Id, req.Code) {
			rsp.Correct = true
		}
	}

	return nil
}
