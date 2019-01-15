package captcha

import (
	"bytes"
	"context"
	"github.com/dchest/captcha"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"teddy-backend/internal/models"
	captchaProto "teddy-backend/internal/proto/captcha"
	"teddy-backend/internal/repositories"
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
	if err := validateGetCaptchaIdReq(req); err != nil {
		return nil, err
	}
	id := captcha.NewLen(int(req.Len))
	return &captchaProto.GetCaptchaIdResp{
		Id: id,
	}, nil
}

func (h *captchaHandler) GetImageData(ctx context.Context, req *captchaProto.GetImageDataReq) (*captchaProto.GetImageDataResp, error) {
	if err := validateGetImageDataReq(req); err != nil {
		return nil, err
	}

	if req.Reload == true {
		if !captcha.Reload(req.Id) {
			return nil, ErrCaptchaNotFount
		}
	}

	imgBuf := &bytes.Buffer{}

	if err := captcha.WriteImage(imgBuf, req.Id, int(req.Width), int(req.Height)); err != nil {
		if err == captcha.ErrNotFound {
			return nil, ErrCaptchaNotFount
		} else {
			log.Error(err)
			return nil, ErrInternal
		}
	}

	return &captchaProto.GetImageDataResp{
		Image: imgBuf.Bytes(),
	}, nil
}

func (h *captchaHandler) GetVoiceData(ctx context.Context, req *captchaProto.GetVoiceDataReq) (*captchaProto.GetVoiceDataResp, error) {
	if err := validateGetVoiceDataReq(req); err != nil {
		return nil, err
	}

	if req.Reload == true {
		if !captcha.Reload(req.Id) {
			return nil, ErrCaptchaNotFount
		}
	}

	voiceBuf := &bytes.Buffer{}

	if err := captcha.WriteAudio(voiceBuf, req.Id, req.Lang); err != nil {
		if err == captcha.ErrNotFound {
			return nil, ErrCaptchaNotFount
		} else {
			log.Error(err)
			return nil, ErrInternal
		}
	}

	return &captchaProto.GetVoiceDataResp{
		VoiceWav: voiceBuf.Bytes(),
	}, nil
}

func (h *captchaHandler) GetRandomById(ctx context.Context, req *captchaProto.GetRandomReq) (*captchaProto.GetRandomResp, error) {
	var resp captchaProto.GetRandomResp
	if err := validateGetRandomReq(req); err != nil {
		return nil, err
	}

	s := ""
	for {
		for i := 0; i < int(req.Len); i++ {
			s += (string)(rand.Intn(10) + 48)
		}

		if _, err := h.repo.FindKeyValuePairByKey(s); err != nil {
			if err == mongo.ErrNoDocuments {
				break
			} else {
				log.Error(err)
				return nil, ErrInternal
			}
		}
	}

	err := h.repo.InsertKeyValuePair(&models.KeyValuePair{
		Key:        req.Id,
		Value:      s,
		ExpireTime: time.Now().Add(Expiration),
	})
	if err != nil {
		log.Error(err)
		return nil, ErrInternal
	}
	resp.Code = s

	return &resp, nil
}

func (h *captchaHandler) Verify(ctx context.Context, req *captchaProto.VerifyReq) (*captchaProto.VerifyResp, error) {
	resp := captchaProto.VerifyResp{
		Correct: false,
	}
	if err := validateVerifyReq(req); err != nil {
		return nil, err
	}

	if req.Type == captchaProto.CaptchaType_RANDOM_BY_ID {
		now := time.Now()
		if _, err := h.repo.FindKeyValuePairByKeyAndValueAndExpire(req.Id, req.Code, now); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, ErrCaptchaNotFount
			} else {
				log.Error(err)
				return nil, ErrInternal
			}
		}
		resp.Correct = true
	} else if req.Type == captchaProto.CaptchaType_IMAGE || req.Type == captchaProto.CaptchaType_VOICE {
		if captcha.VerifyString(req.Id, req.Code) {
			resp.Correct = true
		}
	}

	return &resp, nil
}
