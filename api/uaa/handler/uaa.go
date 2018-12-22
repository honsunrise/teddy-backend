package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lestrrat-go/jwx/jwk"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/api/gin_jwt"
	"github.com/zhsyourai/teddy-backend/common/proto/captcha"
	"github.com/zhsyourai/teddy-backend/common/proto/message"
	"github.com/zhsyourai/teddy-backend/common/proto/uaa"
	"net/http"
	"time"
)

type Uaa struct {
	generator *gin_jwt.JwtGenerator
}

func NewUaaHandler(generator *gin_jwt.JwtGenerator) (*Uaa, error) {
	instance := &Uaa{
		generator: generator,
	}
	return instance, nil
}

func (h *Uaa) HandlerNormal(root gin.IRoutes) {
	root.POST("/register", h.Register)
	root.POST("/login", h.Login)
	root.POST("/changePassword", h.ChangePassword)
	root.POST("/sendEmailCaptcha", h.SendEmailCaptcha)
	root.POST("/resetPassword", h.ResetPassword)
	root.GET("/jwks.json", h.JWKsJSON)
}

func (h *Uaa) HandlerAuth(root gin.IRoutes) {
	root.Any("/logout", h.Logout)
	root.POST("/changePassword", h.ChangePassword)
	root.POST("/resetPassword", h.ResetPassword)
}

func (h *Uaa) HandlerHealth(root gin.IRoutes) {
	root.Any("/", h.ReturnOK)
}

func (h *Uaa) ReturnOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (h *Uaa) Register(ctx *gin.Context) {
	if ctx.Query("type") == "email" {
		// parse body
		type registerReq struct {
			Username string   `json:"username"`
			Password string   `json:"password"`
			Roles    []string `json:"roles"`
			Email    string   `json:"email"`
			Captcha  string   `json:"captcha"`
		}
		var body registerReq
		err := ctx.Bind(&body)
		if err != nil {
			ctx.Error(err)
			return
		}

		// extract the client from the context
		uaaClient, ok := clients.UaaFromContext(ctx)
		if !ok {
			ctx.Error(ErrClientNotFound).SetType(gin.ErrorTypePublic)
			return
		}

		// extract the client from the context
		captchaClient, ok := clients.CaptchaFromContext(ctx)
		if !ok {
			ctx.Error(ErrClientNotFound).SetType(gin.ErrorTypePublic)
			return
		}

		// check email or phone
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if body.Email != "" && body.Captcha != "" {
			rsp, err := captchaClient.Verify(timeoutCtx, &captcha.VerifyReq{
				Type: captcha.CaptchaType_RANDOM_BY_ID,
				Id:   body.Email,
				Code: body.Captcha,
			})
			if err != nil || !rsp.Correct {
				ctx.Error(ErrCaptchaNotCorrect).SetType(gin.ErrorTypePublic)
				return
			}
		} else {
			ctx.Error(ErrCaptchaNotCorrect).SetType(gin.ErrorTypePublic)
			return
		}

		// make request
		timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		response, err := uaaClient.RegisterByNormal(timeoutCtx, &uaa.RegisterNormalReq{
			Username: body.Username,
			Password: body.Password,
			Roles:    body.Roles,
			Contact: &uaa.RegisterNormalReq_Email{
				Email: body.Email,
			},
		})
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"uid":   response.Uid,
			"roles": response.Roles,
		})

		// This step can happen error and will ignore
		messageClient, ok := clients.MessageFromContext(ctx)
		if ok {
			// Send welcome email
			timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			messageClient.SendEmail(timeoutCtx, &message.SendEmailReq{
				Email:    response.Email,
				Topic:    "Welcome " + response.Username,
				Content:  "Hi " + response.Username,
				SendTime: ptypes.TimestampNow(),
			})

			// Send welcome inbox
			timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			messageClient.SendInBox(timeoutCtx, &message.SendInBoxReq{
				Uid:      response.Uid,
				Topic:    "Welcome " + body.Username,
				Content:  "Hi " + body.Username,
				SendTime: ptypes.TimestampNow(),
			})
		}
	} else {
		ctx.Error(ErrRegisterTypeNotSupport).SetType(gin.ErrorTypePublic)
	}
}

func (h *Uaa) Login(ctx *gin.Context) {

	// parse body
	type loginReq struct {
		Principal string `json:"principal"`
		Password  string `json:"password"`
		Captcha   string `json:"captcha"`
	}
	var body loginReq
	err := ctx.Bind(&body)
	if err != nil {
		ctx.Error(err)
		return
	}

	// extract the client from the context
	uaaClient, ok := clients.UaaFromContext(ctx)
	if !ok {
		ctx.Error(ErrClientNotFound).SetType(gin.ErrorTypePublic)
		return
	}

	// make request
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	response, err := uaaClient.VerifyPassword(timeoutCtx, &uaa.VerifyAccountReq{
		Principal: body.Principal,
		Password:  body.Password,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	token, err := h.generator.GenerateJwt(24*time.Hour, 72*time.Hour, jwt.MapClaims{
		"uid":      response.Uid,
		"username": response.Username,
		"roles":    response.Roles,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"type":         "bearer",
	})
}

func (h *Uaa) Logout(ctx *gin.Context) {
	// TODO: may do something
	ctx.Status(http.StatusOK)
}

func (h *Uaa) ChangePassword(ctx *gin.Context) {
	// parse body
	type changePasswordReq struct {
		Principal       string `json:"principal"`
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		CaptchaId       string `json:"captcha_id"`
		CaptchaSolution string `json:"captcha_solution"`
	}
	var body changePasswordReq
	err := ctx.Bind(&body)
	if err != nil {
		ctx.Error(err)
		return
	}

	// extract the client from the context
	captchaClient, ok := clients.CaptchaFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	rsp, err := captchaClient.Verify(timeoutCtx, &captcha.VerifyReq{
		Type: captcha.CaptchaType_IMAGE,
		Id:   body.CaptchaId,
		Code: body.CaptchaSolution,
	})

	if err != nil || !rsp.Correct {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrCaptchaNotCorrect)
		return
	}

	// extract the client from the context
	uaaClient, ok := clients.UaaFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
	}

	// make request
	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = uaaClient.ChangePassword(timeoutCtx, &uaa.ChangePasswordReq{
		Principal:   body.Principal,
		NewPassword: body.NewPassword,
		OldPassword: body.OldPassword,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Status(http.StatusOK)
}

func (h *Uaa) SendEmailCaptcha(ctx *gin.Context) {
	now := time.Now()
	// parse body
	type sendEmailCaptchaReq struct {
		Email           string `json:"email"`
		CaptchaId       string `json:"captcha_id"`
		CaptchaSolution string `json:"captcha_solution"`
	}
	var body sendEmailCaptchaReq
	err := ctx.Bind(&body)
	if err != nil {
		ctx.Error(err)
		return
	}

	// extract the client from the context
	uaaClient, ok := clients.UaaFromContext(ctx)
	if !ok {
		ctx.Error(ErrClientNotFound).SetType(gin.ErrorTypePublic)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tmpAccount, err := uaaClient.GetOne(timeoutCtx, &uaa.GetOneReq{
		Principal: body.Email,
	})

	if err == nil && tmpAccount != nil {
		ctx.Error(ErrAccountExists).SetType(gin.ErrorTypePublic)
		return
	}

	// extract the client from the context
	captchaClient, ok := clients.CaptchaFromContext(ctx)
	if !ok {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	rsp, err := captchaClient.Verify(timeoutCtx, &captcha.VerifyReq{
		Type: captcha.CaptchaType_IMAGE,
		Id:   body.CaptchaId,
		Code: body.CaptchaSolution,
	})

	if err != nil || !rsp.Correct {
		log.Error(ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusBadRequest, ErrCaptchaNotCorrect)
		return
	}

	// extract the client from the context
	messageClient, ok := clients.MessageFromContext(ctx)
	if !ok {
		log.Error("message client not found")
		ctx.AbortWithError(http.StatusInternalServerError, ErrClientNotFound)
		return
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	random, err := captchaClient.GetRandomById(timeoutCtx, &captcha.GetRandomReq{
		Len: 6,
		Id:  body.Email,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Send captcha email
	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = messageClient.SendEmail(ctx, &message.SendEmailReq{
		Email:   body.Email,
		Topic:   "Verify captcha",
		Content: "Your captcha:" + random.Code,
		SendTime: &timestamp.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Uaa) SendPhoneCaptcha(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (h *Uaa) ResetPassword(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (h *Uaa) JWKsJSON(ctx *gin.Context) {
	privKey, err := h.generator.GetJwtPublishKey()
	if err != nil {
		ctx.Error(err)
		return
	}
	key, err := jwk.New(privKey)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, key)
}
