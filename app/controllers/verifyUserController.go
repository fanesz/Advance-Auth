package controller

import (
	"advanceauth/backend/app/handler"
	"advanceauth/backend/app/models"
	"advanceauth/backend/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type VerifyUserController struct {
	Db *gorm.DB
}

func (a *VerifyUserController) ResendVerify(c *gin.Context) {
	// middleware auth
	// getting user uuid from middleware
	UUID := c.GetString("UUID")

	// check is user already verified
	queryUser := a.Db.Model(models.User{UUID: UUID})
	val, _ := handler.QueryValidator(queryUser, c, false)
	if !val {
		return
	}

	// if user already verified
	dataUser := models.User{}
	queryUser.First(&dataUser)
	if dataUser.IsVerified {
		handler.Error(c, http.StatusForbidden, utils.VERIFY_ALREADY_VERIFIED, "User already verified")
		return
	}

	// getting verify token
	verifyUser := models.VerifyUser{
		UUID: UUID,
	}
	queryGetVerifyToken := a.Db.Model(verifyUser)
	val, count := handler.QueryValidator(queryGetVerifyToken, c, true)
	if !val && count == -1 {
		return
	}

	// if no verify token but user not verified (when success register but failed to create verify record)
	if count == 0 {
		// generate verify token
		newVerifyToken := utils.GenerateRandomString(50, true, false)
		verifyUser.VerifyToken = newVerifyToken
		verifyUser.EmailSent = 1
		queryCreateVerifyToken := a.Db.Create(&verifyUser)
		val, _ = handler.QueryValidator(queryCreateVerifyToken, c, false)
		if !val {
			return
		}

		// send email verification
		sendVerifyEmail(dataUser, newVerifyToken)
		handler.Success(c, http.StatusOK, "Email sent", nil)
		return
	}

	// if user not verified
	queryGetVerifyToken.First(&verifyUser)
	diff := utils.GetTimeDifference(verifyUser.UpdatedAt)
	if verifyUser.EmailSent == 3 && diff < 1440 { // if email sent already reach max limit (3)
		handler.Error(c, http.StatusForbidden, utils.VERIFY_MAX_LIMIT_REQUEST, "Max limit request")
		return
	}

	// if email sent already reach max limit (3) but more than 24 hours
	if verifyUser.EmailSent == 3 && diff >= 1440 {
		verifyUser.EmailSent = 1
		queryUpdateVerifyToken := a.Db.Model(&models.VerifyUser{ID: verifyUser.ID}).Updates(verifyUser)
		val, _ = handler.QueryValidator(queryUpdateVerifyToken, c, false)
		if !val {
			return
		}
	}

	// send email verification
	sendVerifyEmail(dataUser, verifyUser.VerifyToken)
	handler.Success(c, http.StatusOK, "Email sent", nil)
}

func (a *VerifyUserController) ValidateVerifyToken(c *gin.Context) {
	// params { token }
	// get params
	verifyToken := c.Param("verify_token")

	// check is verify token exist
	verifyUser := models.VerifyUser{
		VerifyToken: verifyToken,
	}
	queryGetVerifyToken := a.Db.Model(verifyUser).Where("verify_token = ?", verifyToken)
	val, count := handler.QueryValidator(queryGetVerifyToken, c, true)
	if !val && count == -1 {
		return
	}

	// if verify token not exist
	if count == 0 {
		handler.Error(c, http.StatusForbidden, utils.VERIFY_INVALID_TOKEN, "Forbidden")
		return
	}

	// if verify token exist
	queryGetVerifyToken.First(&verifyUser)

	// get user data
	dataUser := models.User{
		UUID: verifyUser.UUID,
	}
	queryGetUser := a.Db.Model(dataUser)
	val, _ = handler.QueryValidator(queryGetUser, c, true)
	if !val {
		return
	}

	// update user is_verified
	queryGetUser.First(&dataUser)
	dataUser.IsVerified = true
	queryUpdateUser := a.Db.Save(&dataUser)
	val, _ = handler.QueryValidator(queryUpdateUser, c, false)
	if !val {
		return
	}

	// delete verify token
	queryDeleteVerifyToken := a.Db.Delete(&verifyUser)
	val, _ = handler.QueryValidator(queryDeleteVerifyToken, c, false)
	if !val {
		return
	}

	handler.Success(c, http.StatusOK, "Account verified", nil)
}

func sendVerifyEmail(dataUser models.User, token string) {
	handler.SendMail(handler.MailInfo{
		EmailTarget: dataUser.Email,
		Subject:     "Account Verification",
		Body: fmt.Sprintf(
			handler.VerifyAccount,
			dataUser.Username,
			utils.GetEnv("FE_URL")+"/user/verify?"+token,
		),
	})
}
