package code_error

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type Error struct {
	traceID  interface{}
	code     int
	httpCode int
	err      string
}

func (e *Error) Error() string {
	return e.err
}

func (e *Error) SetTraceID(traceID interface{}) {
	e.traceID = traceID
}

func DefineCodeError(httpCode, code int, text string) *Error {
	return &Error{
		code:     code,
		err:      text,
		httpCode: httpCode,
	}
}

type Config struct {
	DefaultHttpCode int
	DefaultErrCode  int
}

func NewCodeError(config Config) *codeError {
	if config.DefaultHttpCode == 0 {
		panic(errors.New("please set default http code"))
	}
	if config.DefaultErrCode == 0 {
		panic(errors.New("please set default error code"))
	}

	return &codeError{Config: config}
}

type codeError struct {
	Config
}

func (ne *codeError) AbortWithErrorJSON(ctx *gin.Context, err error) {
	code := ne.DefaultErrCode
	httpCode := ne.DefaultHttpCode
	var traceID interface{} = nil
	switch err.(type) {
	case *Error:
		httpCode = err.(*Error).httpCode
		code = err.(*Error).code
		if err.(*Error).traceID != nil {
			traceID = err.(*Error).traceID
		}
	}

	if traceID != nil {
		ctx.AbortWithStatusJSON(httpCode, gin.H{
			"code":    code,
			"detail":  err.Error(),
			"traceID": traceID,
		})
	} else {
		ctx.AbortWithStatusJSON(httpCode, gin.H{
			"code":   code,
			"detail": err.Error(),
		})
	}
}
