package handler

import "errors"

var ErrClientNotFound = errors.New("client not found")
var ErrCaptchaNotCorrect = errors.New("captcha is not match")

var ErrOrderNotCorrect = errors.New("order not correct must be asc or desc")
var ErrAccountExist = errors.New("account exist")
