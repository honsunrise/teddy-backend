package server

import (
	"errors"
	"github.com/zhsyourai/teddy-backend/common/proto"
)

var ErrPasswordEmpty = errors.New("password can't be empty")
var ErrUsernameEmpty = errors.New("username can't be empty")
var ErrRolesEmpty = errors.New("role can't be empty")
var ErrEmailOrPhoneEmpty = errors.New("email or phone can't be empty")
var ErrEmailAndPhoneExist = errors.New("both email and phone exist")
var ErrOldPasswordEmpty = errors.New("old password empty")
var ErrNewPasswordEmpty = errors.New("new password empty")

func validateRegisterReq(req *proto.RegisterReq) error {
	if req.Password == "" {
		return ErrPasswordEmpty
	} else if req.Username == "" {
		return ErrUsernameEmpty
	} else if len(req.Roles) == 0 {
		return ErrRolesEmpty
	} else {
		if req.Email == "" {
			if req.Phone == "" {
				return ErrEmailOrPhoneEmpty
			}
		} else if req.Phone != "" {
			return ErrEmailAndPhoneExist
		}
	}
	return nil
}

func validateGetByUsernameReq(req *proto.GetByUsernameReq) error {
	if req.Username == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func validateDeleteByUsernameReq(req *proto.DeleteByUsernameReq) error {
	if req.Username == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func validateVerifyPasswordReq(req *proto.VerifyPasswordReq) error {
	if req.Username == "" {
		return ErrUsernameEmpty
	} else if req.Password == "" {
		return ErrPasswordEmpty
	}
	return nil
}

func validateChangePasswordReq(req *proto.ChangePasswordReq) error {
	if req.Username == "" {
		return ErrUsernameEmpty
	} else if req.OldPassword == "" {
		return ErrOldPasswordEmpty
	} else if req.NewPassword == "" {
		return ErrNewPasswordEmpty
	}
	return nil
}
