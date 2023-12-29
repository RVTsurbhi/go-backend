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

// func NewHttpError(description, metadata string, statusCode int) Http {
//     return Http{
//         Description: description,
//         Metadata:    metadata,
//         StatusCode:  statusCode,
//     }
// }

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
		// for _, err := range c.Errors {
		// 	switch e := err.Err.(type) {
		// 	case *CustomError:
		// 		c.AbortWithStatusJSON(e.Code, e)
		// 	default:
		// 		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
		// 	}

		// 	c.Abort()
		// }

		err := c.Errors.Last()
		if err != nil {
			// Tangani error di sini
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

// func ErrorHandler() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Next()
// 		// err := c.Errors.Last()
// 		// if err != nil {
// 		// log.Println("found an err", err)
// 		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
// 		// }

// 		for _, err := range c.Errors {

// 			switch err.Err {
// 			case ErrNotFound:
// 				// c.JSON(-1, gin.H{"error": ErrNotFound.Error()})
// 				c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Not Found"})
// 			case ErrUnauthorized:
// 				c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
// 			// case 'ErrForbidden':
// 			// 	c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Forbidden"})
// 			// case 'ErrBadRequest':
// 			// 	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Bad Request"})
// 			// case 'ErrInternalServer':
// 			// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal Server Error"})
// 			default:
// 				log.Println("found an err2", err.Err)
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
// 				// c.JSON(500, gin.H{
// 				// 	"error": err.Error(),
// 				// })
// 			}
// 			c.Abort()
// 			// etc...
// 		}
// 		// c.JSON(http.StatusInternalServerError, "")
// 	}
// }
