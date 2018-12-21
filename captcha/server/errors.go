package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInternal = status.Error(codes.Internal, "internal")
var ErrCaptchaNotFount = status.Error(codes.NotFound, "captcha not found")
