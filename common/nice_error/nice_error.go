package nice_error

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewNiceError() gin.HandlerFunc {
	ne := niceError{}
	return func(context *gin.Context) {
		ne.process(context)
	}
}

func DefineNiceError(code int, title string, detail string) *NiceError {
	return &NiceError{
		Code:   code,
		Title:  title,
		Detail: detail,
	}
}

type NiceError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (ne *NiceError) Error() string {
	return fmt.Sprintf("<%s> %s", ne.Title, ne.Detail)
}

type niceError struct{}

func (ne *niceError) process(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err != nil {
		realErr := err.Err
		if err.IsType(gin.ErrorTypePublic) {
			switch realErr.(type) {
			case *NiceError:
				toNe := realErr.(*NiceError)
				c.JSON(toNe.Code, gin.H{
					"title":  toNe.Title,
					"detail": toNe.Detail,
				})
			default:
				toErr := realErr.(error)
				c.JSON(http.StatusInternalServerError, gin.H{
					"title":  "Internal Server Error",
					"detail": toErr.Error(),
				})
			}
		} else if err.IsType(gin.ErrorTypeBind) {
			c.JSON(http.StatusBadRequest, gin.H{
				"title":  "Binding Request Error",
				"detail": "request content missing field",
			})
		} else {
			toErr := realErr.(error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Internal Server Error",
				"detail": toErr.Error(),
			})
		}
	}
}
