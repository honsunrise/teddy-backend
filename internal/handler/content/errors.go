package content

import "errors"

var ErrCaptchaNotCorrect = errors.New("captcha is not match")

var ErrOrderNotCorrect = errors.New("order not correct must be asc or desc")
var ErrTagNotCorrect = errors.New("tag not correct must have type and tag")
