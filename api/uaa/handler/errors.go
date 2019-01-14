package handler

import (
	"net/http"
	"teddy-backend/common/nice_error"
)

var ErrCaptchaNotCorrect = nice_error.DefineNiceError(http.StatusBadRequest,
	"correct not correct, please check again")

var ErrRegisterTypeNotSupport = nice_error.DefineNiceError(http.StatusNotFound,
	"register type mistake, please check your request")

var ErrAccountExists = nice_error.DefineNiceError(http.StatusBadRequest,
	"account has been register, please check your request")
