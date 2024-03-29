package gin_jwt

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"time"
)

type GeneratorConfig struct {
	Issuer string

	SigningAlgorithm jose.SignatureAlgorithm

	KeyFunc func() interface{}

	NowFunc func() time.Time
}

type JwtGenerator struct {
	config GeneratorConfig
	jwks   []byte
	signer jose.Signer
}

func NewGinJwtGenerator(config GeneratorConfig) (*JwtGenerator, error) {
	if config.SigningAlgorithm == "" {
		return nil, ErrMissingSigningAlgorithm
	}

	if config.KeyFunc == nil {
		return nil, ErrMissingKeyFunction
	}

	if config.NowFunc == nil {
		config.NowFunc = time.Now
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: config.SigningAlgorithm, Key: config.KeyFunc()}, nil)
	if err != nil {
		return nil, err
	}

	jwk := jose.JSONWebKey{
		Key: config.KeyFunc(),
	}

	thumbprint, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	jwk.KeyID = base64.URLEncoding.EncodeToString(thumbprint)

	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk.Public()},
	}

	jwksResult, err := json.Marshal(jwks)
	if err != nil {
		return nil, err
	}

	return &JwtGenerator{
		config: config,
		jwks:   jwksResult,
		signer: signer,
	}, nil
}

func (g *JwtGenerator) GetJwks() []byte {
	return g.jwks
}

func (g *JwtGenerator) GenerateJwt(timeout time.Duration, subject string, audience []string,
	claims map[string]interface{}) (string, error) {
	now := g.config.NowFunc()
	expire := now.Add(timeout)
	tokenString, err := jwt.Signed(g.signer).
		Claims(jwt.Claims{
			ID:        uuid.Must(uuid.NewRandom()).String(),
			Subject:   subject,
			Issuer:    g.config.Issuer,
			Audience:  audience,
			Expiry:    jwt.NewNumericDate(expire),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		}).
		Claims(claims).
		CompactSerialize()
	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}
