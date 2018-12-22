package handler

import "github.com/zhsyourai/teddy-backend/api/nice_error"

var ErrClientNotFound = nice_error.DefineNiceError(404, "Client not found",
	"server client instance create error, please try again later")

var ErrCaptchaNotCorrect = nice_error.DefineNiceError(404, "Captcha not correct",
	"correct not correct, please check again")

var ErrRegisterTypeNotSupport = nice_error.DefineNiceError(404, "Register type not support",
	"register type mistake, please check your request")

var ErrAccountExists = nice_error.DefineNiceError(404, "Account exists",
	"account has been register, please check your request")
