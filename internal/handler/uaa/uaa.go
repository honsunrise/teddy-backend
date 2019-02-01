package uaa

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"teddy-backend/internal/clients"
	"teddy-backend/internal/gin_jwt"
	"teddy-backend/internal/handler/errors"
	"teddy-backend/internal/proto/captcha"
	"teddy-backend/internal/proto/message"
	"teddy-backend/internal/proto/uaa"
	"time"
)

type Uaa struct {
	generator *gin_jwt.JwtGenerator
	middle    *gin_jwt.JwtMiddleware
}

func NewUaaHandler(middle *gin_jwt.JwtMiddleware, generator *gin_jwt.JwtGenerator) (*Uaa, error) {
	instance := &Uaa{
		generator: generator,
		middle:    middle,
	}
	return instance, nil
}

func (h *Uaa) HandlerNormal(root gin.IRoutes) {
	root.POST("/register", h.Register)
	root.POST("/login", h.Login)
	root.POST("/sendEmailCaptcha", h.SendEmailCaptcha)
	root.POST("/resetPassword", h.ResetPassword)
	root.GET("/jwks.json", h.JWKsJSON)
}

func (h *Uaa) HandlerAuth(root gin.IRoutes) {
	root.POST("/logout", h.Logout)
	root.POST("/changePassword", h.ChangePassword)
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
	uaaClient := clients.UaaFromContext(ctx)
	captchaClient := clients.CaptchaFromContext(ctx)
	messageClient := clients.MessageFromContext(ctx)

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
			errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
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
				errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
				return
			}
		} else {
			errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
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
			errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
			return
		}

		h.middle.AddUser(response.Uid)

		ctx.JSON(http.StatusOK, gin.H{
			"uid":   response.Uid,
			"roles": response.Roles,
		})

		// This step can happen error and will ignore
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
	} else {
		errors.AbortWithErrorJSON(ctx, errors.ErrRegisterTypeNotSupport)
		return
	}
}

func (h *Uaa) Login(ctx *gin.Context) {
	uaaClient := clients.UaaFromContext(ctx)

	// parse body
	type loginReq struct {
		Principal string `json:"principal"`
		Password  string `json:"password"`
		Captcha   string `json:"captcha"`
	}
	var body loginReq
	err := ctx.Bind(&body)
	if err != nil {
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
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
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
		return
	}

	token, err := h.generator.GenerateJwt(24*time.Hour, response.Uid, []string{"uaa", "content", "message"}, jwt.MapClaims{
		"username": response.Username,
	})
	if err != nil {
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
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
	captchaClient := clients.CaptchaFromContext(ctx)
	uaaClient := clients.UaaFromContext(ctx)

	principal := h.middle.ExtractSub(ctx)

	// parse body
	type changePasswordReq struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		CaptchaId       string `json:"captcha_id"`
		CaptchaSolution string `json:"captcha_solution"`
	}
	var body changePasswordReq
	err := ctx.Bind(&body)
	if err != nil {
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
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
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
		return
	}

	// make request
	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = uaaClient.ChangePassword(timeoutCtx, &uaa.ChangePasswordReq{
		Principal:   principal,
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
	messageClient := clients.MessageFromContext(ctx)
	captchaClient := clients.CaptchaFromContext(ctx)
	uaaClient := clients.UaaFromContext(ctx)

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
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tmpAccount, err := uaaClient.GetOne(timeoutCtx, &uaa.GetOneReq{
		Principal: body.Email,
	})

	if err == nil && tmpAccount != nil {
		errors.AbortWithErrorJSON(ctx, errors.ErrAccountExists)
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
		errors.AbortWithErrorJSON(ctx, errors.ErrCaptchaNotCorrect)
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
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", h.generator.GetJwks())
}
