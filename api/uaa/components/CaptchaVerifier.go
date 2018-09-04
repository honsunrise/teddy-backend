package components

import (
	"github.com/zhsyourai/teddy-backend/api/uaa/repositories"
	"github.com/zhsyourai/teddy-backend/common/models"
	"math/rand"
	"time"
)

type CaptchaVerifier interface {
	NextNumberCaptcha(key string, exist time.Time) (string, error)
	Exist(key string, value string) error
}

func NewCaptchaVerifier(repo repositories.KeyValuePairRepository) (CaptchaVerifier, error) {
	i := &simpleCaptchaVerifier{
		repo: repo,
	}
	go i.cleanTack()
	return i, nil
}

type simpleCaptchaVerifier struct {
	repo repositories.KeyValuePairRepository
}

func (t *simpleCaptchaVerifier) cleanTack() {
	for {
		<-time.After(time.Second)
		t.repo.DeleteKeyValuePairLT(time.Now())
	}
}

func (t *simpleCaptchaVerifier) Exist(key string, value string) error {
	now := time.Now()
	_, err := t.repo.FindKeyValuePairByKeyAndValueAndExpire(key, value, now)
	if err != nil {
		return err
	}
	return nil
}

func (t *simpleCaptchaVerifier) NextNumberCaptcha(key string, exist time.Time) (string, error) {
	s := ""
	for i := 0; i < 6; i++ {
		s += (string)(rand.Intn(10) + 48)
	}
	err := t.repo.InsertKeyValuePair(&models.KeyValuePair{
		Key:        key,
		Value:      s,
		ExpireTime: exist,
	})
	if err != nil {
		return "", err
	}
	return s, nil
}
