package uaa

import "errors"

var ErrPasswordEmpty = errors.New("password can't be empty")
var ErrUsernameEmpty = errors.New("username can't be empty")
var ErrRolesEmpty = errors.New("role can't be empty")
var ErrEmailOrPhoneEmpty = errors.New("email or phone can't be empty")
var ErrOldPasswordEmpty = errors.New("old password empty")
var ErrNewPasswordEmpty = errors.New("new password empty")

var ErrAccountExist = errors.New("account exist")
var UserNotFoundErr = errors.New("user not found")
var OldPasswordNotCorrectErr = errors.New("old password not correct")
var PasswordModifyErr = errors.New("password modify error")
