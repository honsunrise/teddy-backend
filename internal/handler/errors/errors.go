package errors

import (
	"net/http"
)

var ErrUnknown = DefineCodeError(http.StatusInternalServerError, ErrCodeUnknown,
	"unknown error")

var ErrBadRequest = DefineCodeError(http.StatusBadRequest, ErrCodeBadRequest,
	"bad request")

var ErrUsernameOrPasswordNotCorrect = DefineCodeError(http.StatusUnauthorized, ErrCodeUsernameOrPasswordNotCorrect,
	"username or password not correct")

var ErrGRPCDial = DefineCodeError(http.StatusInternalServerError, ErrCodeInternal,
	"internal error, try again latter")

var ErrForbidden = DefineCodeError(http.StatusForbidden, ErrCodeForbidden,
	"forbidden")

var ErrUnauthorized = DefineCodeError(http.StatusUnauthorized, ErrCodeUnauthorized,
	"unauthorized")

var ErrCaptchaIDNotFound = DefineCodeError(http.StatusNotFound, ErrCodeCaptchaIDNotFound,
	"captcha id not found")

var ErrCaptchaExtNotSupport = DefineCodeError(http.StatusBadRequest, ErrCodeCaptchaExtNotSupport,
	"captcha format not support")

var ErrCaptchaNotCorrect = DefineCodeError(http.StatusBadRequest, ErrCodeCaptchaNotCorrect,
	"captcha not correct, please check again")

var ErrRegisterTypeNotSupport = DefineCodeError(http.StatusNotFound, ErrCodeRegisterTypeNotSupport,
	"register type mistake, please check your request")

var ErrAccountExists = DefineCodeError(http.StatusBadRequest, ErrCodeAccountExists,
	"account has been register, please check your request")
