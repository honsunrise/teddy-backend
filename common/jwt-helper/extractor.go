package jwt_helper

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type ExtractorConfig struct {
	Realm string

	SigningAlgorithm string

	KeyFunc func() interface{}

	TokenLookup string
}

type Extractor struct {
	config ExtractorConfig
	key    interface{}
	priKey interface{}
}

func NewJwtExtractor(config ExtractorConfig) (*Extractor, error) {
	if config.Realm == "" {
		return nil, ErrMissingRealm
	}

	if config.SigningAlgorithm == "" {
		return nil, ErrMissingSigningAlgorithm
	}

	if config.KeyFunc == nil {
		return nil, ErrMissingKeyFunction
	}

	if config.TokenLookup == "" {
		config.TokenLookup = "Bearer"
	}

	var realKey interface{}
	switch config.SigningAlgorithm {
	case "RS256", "RS384", "RS512":
		if pubKey, ok := config.KeyFunc().(*rsa.PublicKey); ok {
			realKey = pubKey
		} else {
			return nil, ErrInvalidKey
		}
	case "EC256", "EC384", "EC512":
		if pubKey, ok := config.KeyFunc().(*ecdsa.PublicKey); ok {
			realKey = pubKey
		} else {
			return nil, ErrInvalidKey
		}
	case "HS256", "HS384", "HS512":
		if key, ok := config.KeyFunc().([]byte); ok {
			realKey = key
		} else {
			return nil, ErrInvalidKey
		}
	default:
		return nil, ErrNotSupportSigningAlgorithm
	}

	return &Extractor{
		config: config,
		key:    realKey,
	}, nil
}

func (m *Extractor) ExtractToken(origin string) (*jwt.Token, error) {
	token, err := m.extractToken(origin)

	if err != nil {
		return nil, err
	}

	err = m.checkToken(token)

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Extractor) extractToken(originToken string) (*jwt.Token, error) {
	var token string
	var err error

	tmpParts := strings.SplitN(originToken, " ", 2)
	if !(len(tmpParts) == 2 && tmpParts[0] == m.config.TokenLookup) {
		return nil, ErrTokenInvalid
	}

	token = tmpParts[1]

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(m.config.SigningAlgorithm) != token.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		return m.key, nil
	})
}

func (m *Extractor) checkToken(token *jwt.Token) error {
	if !token.Valid {
		return ErrTokenInvalid
	}
	return nil
}
