package errors

import (
	"net/http"
)

var ErrUnknown = DefineCodeError(http.StatusInternalServerError, ErrCodeUnknown,
	"unknown error")

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
