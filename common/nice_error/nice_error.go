package nice_error

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewNiceError() gin.HandlerFunc {
	ne := niceError{}
	return func(context *gin.Context) {
		ne.process(context)
	}
}

func DefineNiceError(code int, text string) *CodeGinError {
	return &CodeGinError{
		Error: gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  errors.New(text),
		},
		Code: code,
	}
}

func DefinePrivateNiceError(code int, text string) *CodeGinError {
	return &CodeGinError{
		Error: gin.Error{
			Type: gin.ErrorTypePrivate,
			Err:  errors.New(text),
		},
		Code: code,
	}
}

type CodeGinError struct {
	gin.Error
	Code int
}

type niceError struct{}

func (ne *niceError) process(c *gin.Context) {
	c.Next()

	var err interface{}
	err = c.Errors.Last()
	if err != nil {
		if codeErr, ok := err.(*CodeGinError); ok {
			title := "internal server error"
			if codeErr.Code == http.StatusForbidden {
				title = "forbidden"
			} else if codeErr.Code == http.StatusUnauthorized {
				title = "unauthorized"
			}
			c.JSON(codeErr.Code, gin.H{
				"title":  title,
				"detail": codeErr.Error.Error(),
			})
		} else {
			if err.(*gin.Error).IsType(gin.ErrorTypeBind) {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":    "binding request error",
					"detail": "request content missing field",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"title":  "internal server error",
					"detail": err.(*gin.Error).Error(),
				})
			}
		}
	}
}
