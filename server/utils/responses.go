package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, message string, status int) {
	c.JSON(status, gin.H{
		"is_success": true,
		"statusCode": status,
		"message":    message,
	})
}

func DataResponse(c *gin.Context, status int, message string, items gin.H) {
	c.JSON(status, gin.H{
		"is_success": true,
		"statusCode": status,
		"message":    message,
		"data":       items["data"],
	})
}

func PageResponse(c *gin.Context, status int, total int64, pageNo int, items gin.H) {
	c.JSON(status, gin.H{
		"is_success": true,
		"statusCode": status,
		"skip":       pageNo,
		"total":      total,
		"items":      items["data"],
	})
}

func Failureresponse(c *gin.Context, errorMsg string, status int) {
	c.JSON(status, gin.H{
		"is_success": false,
		"statusCode": status,
		"message":    errorMsg,
	})
}
