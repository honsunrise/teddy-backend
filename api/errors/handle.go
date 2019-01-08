package errors

import "github.com/gin-gonic/gin"

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

func AbortWithErrorJSON(ctx *gin.Context, err *Error) {
	if err.traceID != nil {
		ctx.AbortWithStatusJSON(err.httpCode, gin.H{
			"code":    err.code,
			"detail":  err.Error(),
			"traceID": err.traceID,
		})
	} else {
		ctx.AbortWithStatusJSON(err.httpCode, gin.H{
			"code":   err.code,
			"detail": err.Error(),
		})
	}
}
