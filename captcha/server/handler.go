package server

import (
	"bytes"
	"context"
	"github.com/dchest/captcha"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/captcha/models"
	"github.com/zhsyourai/teddy-backend/captcha/repositories"
	captchaProto "github.com/zhsyourai/teddy-backend/common/proto/captcha"
	"math/rand"
	"time"
)

const (
	Expiration = 10 * time.Minute
)

func NewCaptchaServer(repo repositories.KeyValuePairRepository) (captchaProto.CaptchaServer, error) {
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

func (h *captchaHandler) GetCaptchaId(ctx context.Context, req *captchaProto.GetCaptchaIdReq) (*captchaProto.GetCaptchaIdResp, error) {
	var resp captchaProto.GetCaptchaIdResp
	if err := validateGetCaptchaIdReq(req); err != nil {
		return nil, err
	}
	id := captcha.NewLen(int(req.Len))
	resp.Id = id
	return &resp, nil
}

func (h *captchaHandler) GetImageData(ctx context.Context, req *captchaProto.GetImageDataReq) (*captchaProto.GetImageDataResp, error) {
	var resp captchaProto.GetImageDataResp
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

func (h *captchaHandler) GetVoiceData(ctx context.Context, req *captchaProto.GetVoiceDataReq) (*captchaProto.GetVoiceDataResp, error) {
	var resp captchaProto.GetVoiceDataResp
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

func (h *captchaHandler) GetRandomById(ctx context.Context, req *captchaProto.GetRandomReq) (*captchaProto.GetRandomResp, error) {
	log.Infof("Get Random number req %v", req)
	var resp captchaProto.GetRandomResp
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

func (h *captchaHandler) Verify(ctx context.Context, req *captchaProto.VerifyReq) (*captchaProto.VerifyResp, error) {
	var resp captchaProto.VerifyResp
	if err := validateVerifyReq(req); err != nil {
		return nil, err
	}
	resp.Correct = false

	if req.Type == captchaProto.CaptchaType_RANDOM_BY_ID {
		now := time.Now()
		_, err := h.repo.FindKeyValuePairByKeyAndValueAndExpire(req.Id, req.Code, now)
		if err != nil {
			return nil, err
		}
		resp.Correct = true
	} else if req.Type == captchaProto.CaptchaType_IMAGE || req.Type == captchaProto.CaptchaType_VOICE {
		if captcha.VerifyString(req.Id, req.Code) {
			resp.Correct = true
		}
	}

	return &resp, nil
}
