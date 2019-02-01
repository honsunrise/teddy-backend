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
	keyFunc  func() interface{}
	nowFunc  func() time.Time
	audience []string
	issuer   string
	subject  string
	id       string
	adapter  persist.Adapter
	enforcer *casbin.SyncedEnforcer
}

func NewGinJwtMiddleware(config MiddlewareConfig, adapter persist.Adapter) (*JwtMiddleware, error) {
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
			if err == ErrForbidden {
				ctx.AbortWithStatus(http.StatusForbidden)
			} else if err == ErrTokenInvalid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}

	if config.ContextKey == "" {
		config.ContextKey = DefaultContextKey
	}

	enforcer, err := casbin.NewSyncedEnforcerSafe(casbin.NewModel(CasbinModel), adapter)
	if err != nil {
		return nil, err
	}

	enforcer.StartAutoLoadPolicy(10 * time.Second)

	return &JwtMiddleware{
		config:   config,
		keyFunc:  config.KeyFunc,
		nowFunc:  config.NowFunc,
		audience: config.Audience,
		issuer:   config.Issuer,
		subject:  config.Subject,
		id:       config.ID,
		adapter:  adapter,
		enforcer: enforcer,
	}, nil
}

func (m *JwtMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := m.extractToken(ctx)
		if err != nil {
			m.config.ErrorHandler(ctx, err)
			return
		}
		sub := ""
		if token != nil {
			sub = token["sub"].(string)
		}
		if !m.enforcer.Enforce(sub, ctx.Request.URL.Path, ctx.Request.Method) {
			m.config.ErrorHandler(ctx, ErrForbidden)
			return
		}

		ctx.Set(m.config.ContextKey, token)
	}
}

func (m *JwtMiddleware) ExtractClaims(ctx *gin.Context, key string) interface{} {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		return token.(map[string]interface{})[key]
	} else {
		// should never happen
		panic(ErrContextNotHaveToken)
	}
}

func (m *JwtMiddleware) ExtractSub(ctx *gin.Context) string {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		if token.(map[string]interface{})["sub"] != nil {
			return token.(map[string]interface{})["sub"].(string)
		}
	}
	return ""
}

func (m *JwtMiddleware) ExtractIss(ctx *gin.Context) string {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		if token.(map[string]interface{})["iss"] != nil {
			return token.(map[string]interface{})["iss"].(string)
		}
	}
	return ""
}

func (m *JwtMiddleware) ExtractJTI(ctx *gin.Context) string {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		if token.(map[string]interface{})["jti"] != nil {
			return token.(map[string]interface{})["jti"].(string)
		}
	}
	return ""
}

func (m *JwtMiddleware) ExtractNBF(ctx *gin.Context) time.Time {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		if token.(map[string]interface{})["nbf"] != nil {
			return time.Unix(int64(token.(map[string]interface{})["nbf"].(float64)), 0)
		}
	}
	return time.Time{}
}

func (m *JwtMiddleware) ExtractEXP(ctx *gin.Context) time.Time {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		if token.(map[string]interface{})["exp"] != nil {
			return time.Unix(int64(token.(map[string]interface{})["exp"].(float64)), 0)
		}
	}
	return time.Time{}
}

func (m *JwtMiddleware) ExtractIAT(ctx *gin.Context) time.Time {
	if token, ok := ctx.Get(m.config.ContextKey); ok {
		if token.(map[string]interface{})["iat"] != nil {
			return time.Unix(int64(token.(map[string]interface{})["iat"].(float64)), 0)
		}
	}
	return time.Time{}
}

func (m *JwtMiddleware) extractToken(ctx *gin.Context) (map[string]interface{}, error) {
	var token string

	parts := strings.SplitN(m.config.TokenLookup, ":", 3)
	switch parts[0] {
	case "header":
		originToken := ctx.Request.Header.Get(parts[1])
		if originToken == "" {
			return nil, nil
		}
		tmpParts := strings.SplitN(originToken, " ", 2)
		if !(len(tmpParts) == 2 && tmpParts[0] == parts[2]) {
			return nil, ErrInvalidAuthHeader
		}
		token = tmpParts[1]
	case "query":
		token = ctx.Query(parts[1])
		if token == "" {
			return nil, nil
		}
	case "cookie":
		token, _ = ctx.Cookie(parts[1])
		if token == "" {
			return nil, nil
		}
	}

	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	c := make(map[string]interface{})
	err = parsedToken.Claims(m.keyFunc(), &c)
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
