package uaa

import (
	"teddy-backend/internal/proto/uaa"
)

func validateRegisterNormalReq(req *uaa.RegisterNormalReq) error {
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

func validateGetOneReq(req *uaa.GetOneReq) error {
	if req.Principal == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func validateUIDReq(req *uaa.UIDReq) error {
	if req.Uid == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func validateVerifyPasswordReq(req *uaa.VerifyAccountReq) error {
	if req.Principal == "" {
		return ErrUsernameEmpty
	} else if req.Password == "" {
		return ErrPasswordEmpty
	}
	return nil
}

func validateChangePasswordReq(req *uaa.ChangePasswordReq) error {
	if req.Principal == "" {
		return ErrUsernameEmpty
	} else if req.OldPassword == "" {
		return ErrOldPasswordEmpty
	} else if req.NewPassword == "" {
		return ErrNewPasswordEmpty
	}
	return nil
}

func validateUpdateSignInReq(req *uaa.UpdateSignInReq) error {
	if req.Ip == "" {

	} else if req.Time == nil {

	}
	return nil
}
