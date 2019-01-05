package nice_error

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func NewNiceError() gin.HandlerFunc {
	ne := niceError{}
	return func(context *gin.Context) {
		ne.process(context)
	}
}

func DefineNiceError(code int, text string) *gin.Error {
	return &gin.Error{
		Type: gin.ErrorTypePublic,
		Err:  errors.New(text),
		Meta: map[string]interface{}{
			"nice_err_code": code,
		},
	}
}

func DefinePrivateNiceError(code int, text string) *gin.Error {
	return &gin.Error{
		Type: gin.ErrorTypePrivate,
		Err:  errors.New(text),
		Meta: map[string]interface{}{
			"nice_err_code": code,
		},
	}
}

type CodeGinError struct {
	Err  gin.Error
	Code int
}

type niceError struct{}

func (ne *niceError) process(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err != nil {
		if err.IsType(gin.ErrorTypeBind) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":    "binding request error",
				"detail": "request content missing field",
			})
		} else if err.IsType(gin.ErrorTypePublic) {
			title := "internal server error"
			code := http.StatusInternalServerError

			if err.Meta != nil {
				value := reflect.ValueOf(err.Meta)
				switch value.Kind() {
				case reflect.Map:
					code = int(reflect.ValueOf(value.MapIndex(reflect.ValueOf("nice_err_code")).Interface()).Int())
				}
			}
			if code == http.StatusForbidden {
				title = "forbidden"
			} else if code == http.StatusUnauthorized {
				title = "unauthorized"
			}
			c.JSON(code, gin.H{
				"title":  title,
				"detail": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "internal server error",
				"detail": err.Error(),
			})
		}
	}
}
