package handler

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/gin-jwt"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/common/errors"
	msgProto "github.com/zhsyourai/teddy-backend/message/proto"
	uaaProto "github.com/zhsyourai/teddy-backend/uaa/proto"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"net/http"
	"time"
)

type Uaa struct {
	enforcer   *casbin.Enforcer
	middleware *gin_jwt.JwtMiddleware
	generator  *gin_jwt.JwtGenerator
}

func NewUaaHandler(enforcer *casbin.Enforcer, middleware *gin_jwt.JwtMiddleware, generator *gin_jwt.JwtGenerator) (*Uaa, error) {
	instance := &Uaa{
		enforcer:   enforcer,
		middleware: middleware,
		generator:  generator,
	}
	return instance, nil
}

func (h *Uaa) Handler(root gin.IRoutes) {
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
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		log.Error("uaa client not found")
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	// check email or phone
	// TODO: Check email or phone

	// make request
	response, err := uaaClient.Register(ctx, &uaaProto.RegisterReq{
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
	messageClient, ok := client.MessageFromContext(ctx)
	if ok {
		// Send welcome email
		messageClient.SendEmail(ctx, &msgProto.SendEmailReq{
			Email:   response.Email,
			Topic:   "Welcome " + response.Username,
			Content: "Hi " + response.Username,
			SendTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		})

		// Send welcome inbox
		messageClient.SendInBox(ctx, &msgProto.SendInBoxReq{
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
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		log.Error("uaa client not found")
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
		return
	}

	// make request
	response, err := uaaClient.VerifyPassword(ctx, &uaaProto.VerifyPasswordReq{
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
	token, err := h.middleware.ExtractToken(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	sub := token.Claims.(jwt.MapClaims)["uid"] // the user that wants to access a resource.
	obj := "uaa.logout"                        // the resource that is going to be accessed.
	act := "read,write"                        // the operation that the user performs on the resource.

	if h.enforcer.Enforce(sub, obj, act) != true {
		ctx.AbortWithStatus(http.StatusForbidden)
	}
	// TODO: may do something
	ctx.Status(http.StatusOK)
}

// Uaa.ChangePassword is called by the API as /uaa/ChangePassword with post body
func (h *Uaa) ChangePassword(ctx *gin.Context) {
	token, err := h.middleware.ExtractToken(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	sub := token.Claims.(jwt.MapClaims)["uid"] // the user that wants to access a resource.
	obj := "uaa.changePassword"                // the resource that is going to be accessed.
	act := "read,write"                        // the operation that the user performs on the resource.

	if h.enforcer.Enforce(sub, obj, act) != true {
		ctx.AbortWithStatus(http.StatusForbidden)
	}

	// parse body
	type changePasswordReq struct {
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		Captcha     string `json:"captcha"`
	}
	var body changePasswordReq

	ctx.Bind(&body)

	// extract the client from the context
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		log.Error("uaa client not found")
		ctx.AbortWithError(http.StatusInternalServerError, errors.ErrClientNotFound)
	}

	// make request
	_, err = uaaClient.ChangePassword(ctx, &uaaProto.ChangePasswordReq{
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

	// parse body
	type sendEmailCaptchaReq struct {
		Email string `json:"email"`
	}
	var body sendEmailCaptchaReq

	ctx.Bind(&body)

	// extract the client from the context

	ctx.Status(http.StatusOK)
}

// Uaa.SendEmailVerify is called by the API as /uaa/sendPhoneCaptcha with post body
func (h *Uaa) SendPhoneCaptcha(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
