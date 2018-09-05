package captcha

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

func NewCaptchaHandler(repo repositories.KeyValuePairRepository) (proto.CaptchaHandler, error) {
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

func (h *captchaHandler) GetImage(ctx context.Context, req *proto.GetImageReq, rsp *proto.GetImageResp) error {
	if err := validateGetImageReq(req); err != nil {
		return err
	}

	id := captcha.New()
	rsp.Id = id
	imgBuf := bytes.NewBuffer(rsp.Image)

	if err := captcha.WriteImage(imgBuf, id, int(req.Width), int(req.Height)); err != nil {
		return err
	}

	return nil
}

func (h *captchaHandler) GetVoice(ctx context.Context, req *proto.GetVoiceReq, rsp *proto.GetVoiceResp) error {
	if err := validateGetVoiceReq(req); err != nil {
		return err
	}

	id := captcha.New()
	rsp.Id = id
	voiceBuf := bytes.NewBuffer(rsp.VoiceWav)

	if err := captcha.WriteAudio(voiceBuf, id, req.Lang); err != nil {
		return err
	}

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
		if captcha.Verify(req.Id, []byte(req.Code)) {
			rsp.Correct = true
		}
	}

	return nil
}
