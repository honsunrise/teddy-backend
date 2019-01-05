package handler

import "github.com/zhsyourai/teddy-backend/common/nice_error"

var ErrClientNotFound = nice_error.DefineNiceError(404, "Client not found", "server client "+
	"instance create error, please try again later")

var ErrCaptchaNotFound = nice_error.DefineNiceError(404, "Captcha not found", "captcha id invalid,"+
	" please check again")
