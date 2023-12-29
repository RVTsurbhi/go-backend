package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func checkUserRole(c *gin.Context, role string) (err error) {
	userRole := c.GetString("role")
	err = nil

	if userRole != role {
		return errors.New("unauthorized")
	}

	return nil
}

func MatchUserType(c *gin.Context, userId string) (err error) {
	userRole := c.GetString("role")
	uid := c.GetString("uid")
	err = nil

	//Note: user can access only his/her profile
	if userRole == "user" && uid != userId {
		return errors.New("unauthorized")
	}
	err = checkUserRole(c, userRole)

	// userIdFromToken := c.MustGet("user_id").(string)

	// if userIdFromToken != userId {
	// 	return errors.New("unauthorized")
	// }

	return nil
}
