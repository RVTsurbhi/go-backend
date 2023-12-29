package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// var ErrNotFound = errors.New("Not Found") // 404
var (
	ErrNotFound     = errors.New("Not Found")
	ErrUnauthorized = errors.New("Unauthorized")
)

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			switch e := err.Err.(type) {
			case *CustomError:
				c.AbortWithStatusJSON(e.Code, e)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
			c.Abort()
		}
	}
}

