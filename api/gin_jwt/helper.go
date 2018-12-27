package gin_jwt

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func _fetch(rawUrl string) (*jose.JSONWebKeySet, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	var src []byte
	switch u.Scheme {
	case "http", "https":
		res, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, err
		}

		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		src = buf
	case "file":
		f, err := os.Open(u.Path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		src = buf
	default:
		return nil, err
	}
	jwks := jose.JSONWebKeySet{}
	err = json.Unmarshal(src, &jwks)
	if err != nil {
		return nil, err
	}
	return &jwks, nil
}

func RemoteFetchFunc(rawUrl string, cacheTimeout time.Duration) func() interface{} {
	var jwks *jose.JSONWebKeySet
	var err error
	expTime := time.Time{}
	return func() interface{} {
		if expTime.Before(time.Now()) || jwks == nil {
			expTime = time.Now().Add(cacheTimeout)
			jwks, err = _fetch(rawUrl)
			if err != nil {
				log.Errorf("fetch remote jwks error: %v", err)
				return nil
			}
		}
		return jwks.Keys[0].Public().Key
	}
}
