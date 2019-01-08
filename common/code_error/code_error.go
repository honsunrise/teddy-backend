package code_error

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
)

func NewCodeError(filterType gin.ErrorType) gin.HandlerFunc {
	ne := codeError{
		filterType: filterType,
	}
	return func(context *gin.Context) {
		ne.process(context)
	}
}

func DefineCodeError(code int, text string) *gin.Error {
	return &gin.Error{
		Type: gin.ErrorTypePublic,
		Err:  errors.New(text),
		Meta: map[string]interface{}{
			"_err_code_": code,
		},
	}
}

type codeError struct {
	filterType gin.ErrorType
}

func (ne *codeError) process(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()

	if err != nil && err.IsType(ne.filterType) && c.Writer.Size() <= 0 {
		code := -1
		if err.Meta != nil {
			value := reflect.ValueOf(err.Meta)
			switch value.Kind() {
			case reflect.Map:
				code = int(reflect.ValueOf(value.MapIndex(reflect.ValueOf("_err_code_")).Interface()).Int())
			}
		}
		c.JSON(-1, gin.H{
			"code":   code,
			"detail": err.Error(),
		})
	}
}
