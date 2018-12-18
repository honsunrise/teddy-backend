package server

import "errors"

var ErrAccountExist = errors.New("account exist")
var UserNotFoundErr = errors.New("user not found")
var OldPasswordNotCorrectErr = errors.New("old password not correct")
var PasswordModifyErr = errors.New("password modify error")
