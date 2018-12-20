package server

import (
	"errors"
	"github.com/zhsyourai/teddy-backend/common/proto"
)

var ErrPasswordEmpty = errors.New("password can't be empty")
var ErrUsernameEmpty = errors.New("username can't be empty")
var ErrRolesEmpty = errors.New("role can't be empty")
var ErrEmailOrPhoneEmpty = errors.New("email or phone can't be empty")
var ErrOldPasswordEmpty = errors.New("old password empty")
var ErrNewPasswordEmpty = errors.New("new password empty")

func validateRegisterNormalReq(req *proto.RegisterNormalReq) error {
	if req.Password == "" {
		return ErrPasswordEmpty
	} else if req.Username == "" {
		return ErrUsernameEmpty
	} else if len(req.Roles) == 0 {
		return ErrRolesEmpty
	} else {
		if req.GetContact() == nil {
			return ErrEmailOrPhoneEmpty
		}
	}
	return nil
}

func validateGetOneReq(req *proto.GetOneReq) error {
	if req.Principal == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func validateUIDReq(req *proto.UIDReq) error {
	if req.Uid == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func validateVerifyPasswordReq(req *proto.VerifyAccountReq) error {
	if req.Principal == "" {
		return ErrUsernameEmpty
	} else if req.Password == "" {
		return ErrPasswordEmpty
	}
	return nil
}

func validateChangePasswordReq(req *proto.ChangePasswordReq) error {
	if req.Principal == "" {
		return ErrUsernameEmpty
	} else if req.OldPassword == "" {
		return ErrOldPasswordEmpty
	} else if req.NewPassword == "" {
		return ErrNewPasswordEmpty
	}
	return nil
}

func validateUpdateSignInReq(req *proto.UpdateSignInReq) error {
	if req.Ip == "" {

	} else if req.Time == nil {

	}
	return nil
}
