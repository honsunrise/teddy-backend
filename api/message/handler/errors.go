package handler

import "errors"

var ErrClientNotFound = errors.New("client not found")
var ErrCaptchaNotCorrect = errors.New("captcha is not match")
