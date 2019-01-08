package errors

import (
	"net/http"
)

var ErrGRPCDial = DefineCodeError(http.StatusInternalServerError, ErrCodeInternal,
	"internal error, try again latter")

var ErrUnknown = DefineCodeError(http.StatusInternalServerError, ErrCodeUnknown,
	"unknown error")

var ErrCaptchaIDNotFound = DefineCodeError(http.StatusNotFound, ErrCodeCaptchaIDNotFound,
	"captcha id not found")

var ErrCaptchaExtNotSupport = DefineCodeError(http.StatusBadRequest, ErrCodeCaptchaExtNotSupport,
	"captcha format not support")
