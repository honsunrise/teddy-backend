package gin_jwt

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
	"time"
)

const DefaultContextKey = "_JWT_TOKEN_KEY_"
const DefaultLeeway = 1.0 * time.Minute

type MiddlewareConfig struct {
	Realm            string
	SigningAlgorithm jose.SignatureAlgorithm
	KeyFunc          func() interface{}
	NowFunc          func() time.Time
	ErrorHandler     func(ctx *gin.Context, err error)
	TokenLookup      string
	ContextKey       string
	Audience         []string
	Issuer           string
	Subject          string
	ID               string
}

type JwtMiddleware struct {
	config   MiddlewareConfig
	key      interface{}
	priKey   interface{}
	nowFunc  func() time.Time
	audience []string
	issuer   string
	subject  string
	id       string
}

func NewGinJwtMiddleware(config MiddlewareConfig) (*JwtMiddleware, error) {
	if config.Realm == "" {
		return nil, ErrMissingRealm
	}

	if config.SigningAlgorithm == "" {
		return nil, ErrMissingSigningAlgorithm
	}

	if config.NowFunc == nil {
		config.NowFunc = time.Now
	}

	if config.KeyFunc == nil {
		return nil, ErrMissingKeyFunction
	}

	if config.TokenLookup == "" {
		config.TokenLookup = "header:Authorization:Bearer"
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = func(ctx *gin.Context, err error) {
			ctx.Header("WWW-Authenticate", "JWT realm="+config.Realm)
			ctx.AbortWithError(http.StatusUnauthorized, err)
		}
	}

	if config.ContextKey == "" {
		config.ContextKey = DefaultContextKey
	}

	return &JwtMiddleware{
		config:   config,
		key:      config.KeyFunc(),
		nowFunc:  config.NowFunc,
		audience: config.Audience,
		issuer:   config.Issuer,
		subject:  config.Subject,
		id:       config.ID,
	}, nil
}

func (m *JwtMiddleware) Handler(isOptional bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := m.extractToken(ctx)

		if !isOptional && err != nil {
			m.config.ErrorHandler(ctx, err)
			return
		}

		if err == nil {
			ctx.Set(m.config.ContextKey, token)
		}

		ctx.Next()
	}
}

func (m *JwtMiddleware) ExtractToken(ctx *gin.Context) (map[string]interface{}, error) {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		return token.(map[string]interface{}), nil
	} else {
		return nil, ErrContextNotHaveToken
	}
}

func (m *JwtMiddleware) extractToken(ctx *gin.Context) (map[string]interface{}, error) {
	var token, originToken string

	parts := strings.SplitN(m.config.TokenLookup, ":", 3)
	switch parts[0] {
	case "header":
		originToken = ctx.Request.Header.Get(parts[1])
		if originToken == "" {
			return nil, ErrEmptyAuthHeader
		}
	case "query":
		originToken = ctx.Query(parts[1])
		if originToken == "" {
			return nil, ErrEmptyQueryToken
		}
	case "cookie":
		originToken, _ = ctx.Cookie(parts[1])
		if originToken == "" {
			return nil, ErrEmptyCookieToken
		}
	}

	tmpParts := strings.SplitN(originToken, " ", 2)
	if !(len(tmpParts) == 2 && tmpParts[0] == parts[2]) {
		return nil, ErrInvalidAuthHeader
	}

	token = tmpParts[1]

	parsedToken, err := jwt.ParseSigned(token)
	c := make(map[string]interface{})
	err = parsedToken.Claims(m.key, &c)
	if err != nil {
		return nil, err
	}

	if m.issuer != "" && m.issuer != c["iss"] {
		return nil, ErrTokenInvalid
	}

	if m.subject != "" && m.subject != c["sub"] {
		return nil, ErrTokenInvalid
	}

	if m.id != "" && m.id != c["jti"] {
		return nil, ErrTokenInvalid
	}

	if len(m.audience) != 0 {
		if tmp, ok := c["aud"].([]interface{}); ok {
			aud := make([]string, len(tmp))
			for i, v := range tmp {
				if aud[i], ok = v.(string); !ok {
					return nil, ErrTokenInvalid
				}
			}
			for _, v := range m.audience {
				find := false
				for _, a := range aud {
					if a == v {
						find = true
					}
				}
				if !find {
					return nil, ErrTokenInvalid
				}
			}
		} else {
			return nil, ErrTokenInvalid
		}
	}

	now := m.nowFunc()
	if nbf, ok := c["nbf"].(int64); !ok || now.Add(DefaultLeeway).Before(time.Unix(nbf, 0)) {
		return nil, ErrTokenInvalid
	}

	if exp, ok := c["exp"].(int64); !ok || now.Add(-DefaultLeeway).After(time.Unix(exp, 0)) {
		return nil, ErrTokenInvalid
	}

	return c, nil
}
