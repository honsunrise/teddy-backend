package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInternal = status.Error(codes.Internal, "internal")
var ErrBehaviorExists = status.Error(codes.AlreadyExists, "behavior exists")
var ErrSegmentExists = status.Error(codes.AlreadyExists, "segment exists")
var ErrSegmentNotExists = status.Error(codes.NotFound, "segment not exists")
var ErrInfoNotExists = status.Error(codes.NotFound, "info not exists")
