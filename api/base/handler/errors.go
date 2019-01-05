package handler

import (
	"github.com/zhsyourai/teddy-backend/common/nice_error"
	"net/http"
)

var ErrCaptchaNotFound = nice_error.DefineNiceError(http.StatusNotFound, "captcha id invalid,"+
	" please check again")
