package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lestrrat-go/jwx/jwk"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/clients"
	"github.com/zhsyourai/teddy-backend/api/gin_jwt"
	"github.com/zhsyourai/teddy-backend/common/errors"
	"github.com/zhsyourai/teddy-backend/common/proto"
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
	type okResp struct {
		Status string `json:"status"`
	}
	var jsonResp okResp
	jsonResp.Status = "OK"
	ctx.JSON(http.StatusOK, &jsonResp)
}

// Uaa.Register is called by the API as /uaa/Register with post body
func (h *Uaa) Register(ctx *gin.Context) {
	// Now time
	now := time.Now()

	// parse body
	type registerReq struct {
		Username string   `json:"username"`
		Password string   `json:"password"`
		Roles    []string `json:"roles"`
		Email    string   `json:"email,omitempty"`
		Phone    string   `json:"phone,omitempty"`
		Captcha  string   `json:"captcha,omitempty"`
	}
	var body registerReq

	ctx.Bind(&body)

	// extract the client from the context
	uaaClient, ok := clients.UaaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	// extract the client from the context
	captchaClient, ok := clients.CaptchaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	// check email or phone
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if body.Email != "" && body.Captcha != "" {
		rsp, err := captchaClient.Verify(timeoutCtx, &proto.VerifyReq{
			Type: proto.CaptchaType_RANDOM_BY_ID,
			Id:   body.Email,
			Code: body.Captcha,
		})
		if err != nil || !rsp.Correct {
			log.Error(errors.ErrCaptchaNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, errors.ErrCaptchaNotCorrect)
			return
		}
	} else if body.Phone != "" && body.Captcha != "" {
		rsp, err := captchaClient.Verify(timeoutCtx, &proto.VerifyReq{
			Type: proto.CaptchaType_RANDOM_BY_ID,
			Id:   body.Phone,
			Code: body.Captcha,
		})
		if err != nil || !rsp.Correct {
			log.Error(errors.ErrCaptchaNotCorrect)
			ctx.AbortWithError(http.StatusInternalServerError, errors.ErrCaptchaNotCorrect)
			return
		}
	} else {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrCaptchaNotCorrect)
		return
	}

	// make request
	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	response, err := uaaClient.Register(timeoutCtx, &proto.RegisterReq{
		Username: body.Username,
		Password: body.Password,
		Roles:    body.Roles,
		Email:    body.Email,
		Phone:    body.Phone,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type registerResp struct {
		Uid        string    `json:"uid"`
		CreateDate time.Time `json:"create_date"`
	}
	var jsonResp registerResp
	jsonResp.Uid = response.Uid
	jsonResp.CreateDate = time.Unix(response.CreateDate.Seconds, int64(response.CreateDate.Nanos))

	ctx.JSON(http.StatusOK, &jsonResp)

	// This step can happen error and will ignore
	messageClient, ok := clients.MessageFromContext(ctx)
	if ok {
		// Send welcome email
		timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		messageClient.SendEmail(timeoutCtx, &proto.SendEmailReq{
			Email:   response.Email,
			Topic:   "Welcome " + response.Username,
			Content: "Hi " + response.Username,
			SendTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		})

		// Send welcome inbox
		timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		messageClient.SendInBox(timeoutCtx, &proto.SendInBoxReq{
			Uid:     response.Uid,
			Topic:   "Welcome " + body.Username,
			Content: "Hi " + body.Username,
			SendTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		})
	}
}

// Uaa.Login is called by the API as /uaa/Login with post body
func (h *Uaa) Login(ctx *gin.Context) {

	// parse body
	type loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Captcha  string `json:"captcha"`
	}
	var body loginReq

	ctx.Bind(&body)

	// extract the client from the context
	uaaClient, ok := clients.UaaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	// make request
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	response, err := uaaClient.VerifyPassword(timeoutCtx, &proto.VerifyPasswordReq{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	token, err := h.generator.GenerateJwt(24*time.Hour, 72*time.Hour, jwt.MapClaims{
		"uid":      response.Uid,
		"username": response.Username,
		"roles":    response.Roles,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type registerResp struct {
		Type        string `json:"type"`
		AccessToken string `json:"access_token"`
	}
	jsonResp := registerResp{
		AccessToken: token,
		Type:        "bearer",
	}

	ctx.JSON(http.StatusOK, &jsonResp)
}

// Uaa.Logout is called by the API as /uaa/Logout
func (h *Uaa) Logout(ctx *gin.Context) {
	// TODO: may do something
	ctx.Status(http.StatusOK)
}

// Uaa.ChangePassword is called by the API as /uaa/ChangePassword with post body
func (h *Uaa) ChangePassword(ctx *gin.Context) {
	// parse body
	type changePasswordReq struct {
		Username        string `json:"username"`
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		CaptchaId       string `json:"captcha_id"`
		CaptchaSolution string `json:"captcha_solution"`
	}
	var body changePasswordReq

	ctx.Bind(&body)

	// extract the client from the context
	captchaClient, ok := clients.CaptchaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	rsp, err := captchaClient.Verify(timeoutCtx, &proto.VerifyReq{
		Type: proto.CaptchaType_RANDOM_BY_ID,
		Id:   body.CaptchaId,
		Code: body.CaptchaSolution,
	})

	if err != nil || !rsp.Correct {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrCaptchaNotCorrect)
		return
	}

	// extract the client from the context
	uaaClient, ok := clients.UaaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
	}

	// make request
	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = uaaClient.ChangePassword(timeoutCtx, &proto.ChangePasswordReq{
		Username:    body.Username,
		NewPassword: body.NewPassword,
		OldPassword: body.OldPassword,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Status(http.StatusOK)
}

// Uaa.SendEmailVerify is called by the API as /uaa/sendEmailCaptcha with post body
func (h *Uaa) SendEmailCaptcha(ctx *gin.Context) {
	now := time.Now()
	// parse body
	type sendEmailCaptchaReq struct {
		Email           string `json:"email"`
		CaptchaId       string `json:"captcha_id"`
		CaptchaSolution string `json:"captcha_solution"`
	}
	var body sendEmailCaptchaReq

	ctx.Bind(&body)

	// extract the client from the context
	captchaClient, ok := clients.CaptchaFromContext(ctx)
	if !ok {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	rsp, err := captchaClient.Verify(timeoutCtx, &proto.VerifyReq{
		Type: proto.CaptchaType_IMAGE,
		Id:   body.CaptchaId,
		Code: body.CaptchaSolution,
	})

	if err != nil || !rsp.Correct {
		log.Error(errors.ErrCaptchaNotCorrect)
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrCaptchaNotCorrect)
		return
	}

	// extract the client from the context
	messageClient, ok := clients.MessageFromContext(ctx)
	if !ok {
		log.Error("message client not found")
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	captcha, err := captchaClient.GetRandomById(timeoutCtx, &proto.GetRandomReq{
		Len: 6,
		Id:  body.Email,
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Send captcha email
	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = messageClient.SendEmail(ctx, &proto.SendEmailReq{
		Email:   body.Email,
		Topic:   "Verify captcha",
		Content: "Your captcha:" + captcha.Code,
		SendTime: &timestamp.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	})
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
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
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	key, err := jwk.New(privKey)
	if err != nil {
		log.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, key)
}
