package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func VerifyUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("You are not authorised to access this!")
		return err
	}
	return err
}

func MatchToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("You are not authorised to access this!Unauthorized to access this resource")
		return err
	}
	err = VerifyUserType(c, userType)
	return err
}
