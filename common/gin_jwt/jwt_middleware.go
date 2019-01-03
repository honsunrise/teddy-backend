package gin_jwt

import (
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
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
	Realm        string
	KeyFunc      func() interface{}
	NowFunc      func() time.Time
	ErrorHandler func(ctx *gin.Context, err error)
	TokenLookup  string
	ContextKey   string
	Audience     []string
	Issuer       string
	Subject      string
	ID           string
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
	adapter  *persist.Adapter
	enforcer *casbin.Enforcer
}

func NewGinJwtMiddleware(config MiddlewareConfig, adapter *persist.Adapter) (*JwtMiddleware, error) {
	if config.Realm == "" {
		return nil, ErrMissingRealm
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

	enforcer, err := casbin.NewEnforcerSafe(casbin.NewModel(CasbinModel), adapter)
	if err != nil {
		return nil, err
	}

	return &JwtMiddleware{
		config:   config,
		key:      config.KeyFunc(),
		nowFunc:  config.NowFunc,
		audience: config.Audience,
		issuer:   config.Issuer,
		subject:  config.Subject,
		id:       config.ID,
		adapter:  adapter,
		enforcer: enforcer,
	}, nil
}

func (m *JwtMiddleware) Handler(isOptional bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := m.extractToken(ctx)

		if !isOptional {
			if err != nil {
				m.config.ErrorHandler(ctx, err)
				return
			}
			user := token["sub"].(string)
			if !m.enforcer.Enforce(user, ctx.Request.URL, ctx.Request.Method) {
				m.config.ErrorHandler(ctx, ErrForbidden)
			}
		}

		if err == nil {
			ctx.Set(m.config.ContextKey, token)
		}

		ctx.Next()
	}
}

func (m *JwtMiddleware) ExtractClaims(ctx *gin.Context, key string) (interface{}, error) {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		return token.(map[string]interface{})[key], nil
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
	if err != nil {
		return nil, ErrTokenInvalid
	}
	c := make(map[string]interface{})
	err = parsedToken.Claims(m.key, &c)
	if err != nil {
		if err == jose.ErrUnsupportedKeyType {
			return nil, ErrInvalidKey
		}
		return nil, ErrTokenInvalid
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
						break
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
	if nbf, ok := c["nbf"].(float64); !ok || now.Add(DefaultLeeway).Before(time.Unix(int64(nbf), 0)) {
		return nil, ErrTokenInvalid
	}

	if exp, ok := c["exp"].(float64); !ok || now.Add(-DefaultLeeway).After(time.Unix(int64(exp), 0)) {
		return nil, ErrTokenInvalid
	}

	return c, nil
}
