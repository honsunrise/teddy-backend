package handler

import (
	"context"
	"encoding/json"
	"github.com/casbin/casbin"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
	"time"
)

type Uaa struct {
	enforcer *casbin.Enforcer
}

func NewUaaHandler(enforcer *casbin.Enforcer) (*Uaa, error) {
	instance := &Uaa{
		enforcer: enforcer,
	}
	return instance, nil
}

// Uaa.Register is called by the API as /uaa/Register with post body
func (e *Uaa) Register(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Uaa.Register request")

	// check method
	if req.Method != "POST" {
		return errors.BadRequest("com.micro.api.uaa", "require post")
	}

	// let's make sure we get json
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.uaa", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.uaa", "expect application/json")
	}

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

	json.Unmarshal([]byte(req.Body), &body)

	// extract the client from the context
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		log.Error("uaa client not found")
		return errors.InternalServerError("com.teddy.api.uaa.register", "uaa client not found")
	}

	// check email or phone
	// TODO: Check email or phone

	// make request
	response, err := uaaClient.Register(ctx, &proto.RegisterReq{
		Username: body.Username,
		Password: body.Password,
		Roles:    body.Roles,
		Email:    body.Email,
		Phone:    body.Phone,
	})
	if err != nil {
		log.Error(err)
		return errors.InternalServerError("com.teddy.api.uaa.register", err.Error())
	}

	type registerResp struct {
		Uid        string    `json:"uid"`
		CreateDate time.Time `json:"create_date"`
	}
	var jsonResp registerResp
	jsonResp.Uid = response.Uid
	jsonResp.CreateDate = time.Unix(response.CreateDate.Seconds, int64(response.CreateDate.Nanos))

	b, _ := json.Marshal(jsonResp)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}

// Uaa.Login is called by the API as /uaa/Login with post body
func (e *Uaa) Login(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Uaa.Login request")

	// check method
	if req.Method != "POST" {
		return errors.BadRequest("com.micro.api.uaa", "require post")
	}

	// let's make sure we get json
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.uaa", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.uaa", "expect application/json")
	}

	// parse body
	type loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Captcha  string `json:"captcha"`
	}
	var body loginReq

	json.Unmarshal([]byte(req.Body), &body)

	// extract the client from the context
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		log.Error("uaa client not found")
		return errors.InternalServerError("com.teddy.api.uaa.login", "uaa client not found")
	}

	// make request
	response, err := uaaClient.VerifyPassword(ctx, &proto.VerifyPasswordReq{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		log.Error(err)
		return errors.InternalServerError("com.teddy.api.uaa.login", err.Error())
	}

	type registerResp struct {
		Uid string `json:"uid"`
	}
	var jsonResp registerResp
	jsonResp.Uid = response.Uid

	b, _ := json.Marshal(jsonResp)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}

// Uaa.Logout is called by the API as /uaa/Logout with post body
func (e *Uaa) Logout(context.Context, *api.Request, *api.Response) error {
	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.
	if e.enforcer.Enforce(sub, obj, act) != true {

	}
	return nil
}

// Uaa.ChangePassword is called by the API as /uaa/ChangePassword with post body
func (e *Uaa) ChangePassword(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Uaa.ChangePassword request")

	// check method
	if req.Method != "POST" {
		return errors.BadRequest("com.micro.api.uaa", "require post")
	}

	// let's make sure we get json
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.uaa", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.uaa", "expect application/json")
	}

	// parse body
	type changePasswordReq struct {
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		Captcha     string `json:"captcha"`
	}
	var body changePasswordReq

	json.Unmarshal([]byte(req.Body), &body)

	// extract the client from the context
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		log.Error("uaa client not found")
		return errors.InternalServerError("com.teddy.api.uaa.changePassword", "uaa client not found")
	}

	// make request
	_, err := uaaClient.ChangePassword(ctx, &proto.ChangePasswordReq{
		Username:    body.Username,
		NewPassword: body.NewPassword,
		OldPassword: body.OldPassword,
	})
	if err != nil {
		log.Error(err)
		return errors.InternalServerError("com.teddy.api.uaa.changePassword", err.Error())
	}

	rsp.StatusCode = 200

	return nil
}

// Uaa.SendEmailVerify is called by the API as /uaa/sendEmailCaptcha with post body
func (e *Uaa) SendEmailCaptcha(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Uaa.SendEmailCaptcha request")

	// check method
	if req.Method != "POST" {
		return errors.BadRequest("com.micro.api.uaa", "require post")
	}

	// let's make sure we get json
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.uaa", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.uaa", "expect application/json")
	}

	// parse body
	type changePasswordReq struct {
		Email string `json:"email"`
	}
	var body changePasswordReq

	json.Unmarshal([]byte(req.Body), &body)

	// extract the client from the context

	rsp.StatusCode = 200

	return nil
}

// Uaa.SendEmailVerify is called by the API as /uaa/sendPhoneCaptcha with post body
func (e *Uaa) SendPhoneCaptcha(ctx context.Context, req *api.Request, rsp *api.Response) error {
	panic("implement me")
}
