package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"advanceauth/backend/app/handler"
	"advanceauth/backend/app/models"
	"advanceauth/backend/app/utils"
	"net/http"
	"strconv"
)

type ResetPwController struct {
	Db *gorm.DB
}

func (a *ResetPwController) ResetRequest(c *gin.Context) {
	// body { email, ip_address, device }
	// binding request
	var bodyReq interface{}
	if err := handler.BindAndStructValidator(c, &bodyReq, false); err != nil {
		return
	}
	bodyReqMap := bodyReq.(map[string]interface{})
	reqUser := models.User{
		Email: bodyReqMap["email"].(string),
	}
	reqReset := models.ResetPassword{
		IPAddress: bodyReqMap["ip_address"].(string),
		Device:    bodyReqMap["device"].(string),
	}

	// check is user exist
	queryUser := a.Db.Model(reqUser).Where("email = ?", reqUser.Email)
	val, count := handler.QueryValidator(queryUser, c, true)
	if !val && count == -1 {
		return
	}

	// if user not exist
	if count == 0 {
		handler.Success(c, http.StatusOK, "Request Success", nil)
		return
	}

	// if user exist
	queryUser.First(&reqUser)
	// reqReset.UUID = reqUser.UUID

	// check is user already request
	resetpwData := models.ResetPassword{
		UUID: reqReset.UUID,
	}

	queryResetPassword := a.Db.Model(resetpwData)
	val, count = handler.QueryValidator(queryResetPassword, c, true)
	if !val && count == -1 {
		return
	}

	// if user already request
	if count != 0 {
		queryResetPassword.First(&resetpwData)
		diff := utils.GetTimeDifference(resetpwData.UpdatedAt)
		resetpwData.IPAddress = reqReset.IPAddress
		resetpwData.Device = reqReset.Device

		// check is user already request 3 times
		if resetpwData.EmailSent >= 3 {
			// reset stored token if previous requested token more than 10 minutes
			if diff > 10 {
				resetpwData.ResetToken = utils.GenerateRandomString(50, true, false)
				resetpwData.EmailSent = 1
				queryUpdate := a.Db.Model(&models.ResetPassword{ID: resetpwData.ID}).Updates(resetpwData)
				val, _ = handler.QueryValidator(queryUpdate, c, false)
				if !val {
					return
				}
				sendResetTokenEmail(reqUser, resetpwData, resetpwData.ResetToken)
				handler.Success(c, http.StatusOK, "New reset token has been sent", nil)
				return
			}
			handler.Error(c, http.StatusForbidden, utils.RESETPW_MAX_LIMIT_REQUEST, strconv.Itoa(10-diff))
			return
		}

		// resend reset token
		resetpwData.EmailSent += 1
		queryUpdate := a.Db.Model(&models.ResetPassword{ID: resetpwData.ID}).Updates(resetpwData)
		val, _ = handler.QueryValidator(queryUpdate, c, false)
		if !val {
			return
		}
		sendResetTokenEmail(reqUser, resetpwData, resetpwData.ResetToken)
		handler.Success(c, http.StatusOK, "Reset token has been resent", nil)
		return
	}

	// if user not request yet
	if count == 0 {
		reqReset.ResetToken = utils.GenerateRandomString(50, false, true)
		reqReset.EmailSent = 1
		quaryCreate := a.Db.Create(&reqReset)
		val, _ = handler.QueryValidator(quaryCreate, c, false)
		if !val {
			return
		}
		sendResetTokenEmail(reqUser, reqReset, reqReset.ResetToken)
	}

	handler.Success(c, http.StatusOK, "Request Success", nil)
}

func (a *ResetPwController) ValidateResetRequest(c *gin.Context) {
	// params { token }, body { ip_address, device }
	// get params and binding request
	reqResetpw := models.ResetPassword{}
	if err := handler.BindAndStructValidator(c, &reqResetpw, false); err != nil {
		return
	}
	reqResetpw.ResetToken = c.Query("reset_token")

	// check is reset token exist
	queryResetPassword := a.Db.Model(reqResetpw).Where("reset_token = ?", reqResetpw.ResetToken)
	val, count := handler.QueryValidator(queryResetPassword, c, true)
	if !val && count == -1 {
		return
	}

	// if reset token not exist
	if count == 0 {
		handler.Error(c, http.StatusForbidden, utils.RESETPW_INVALID_TOKEN, "Invalid reset token")
		return
	}

	// if reset token exist
	resetpwData := models.ResetPassword{}
	queryResetPassword.First(&resetpwData)
	diff := utils.GetTimeDifference(resetpwData.UpdatedAt)

	// check is reset token expired
	if diff > 10 {
		handler.Error(c, http.StatusForbidden, utils.RESETPW_EXPIRED_TOKEN, "Reset token expired")
		return
	}

	// if reset token not expired
	// check is ip address and device same
	if resetpwData.IPAddress != reqResetpw.IPAddress || resetpwData.Device != reqResetpw.Device {
		handler.Error(c, http.StatusForbidden, utils.AUTH_DIFFERENT_IP_OR_DEVICE, "Invalid IPAddress or Device")
		return
	}

	// if ip address and device same
	handler.Success(c, http.StatusOK, "Token, Device, and IP Address Confirmed", nil)
}

func (a *ResetPwController) ResetPassword(c *gin.Context) {
	// params { reset_token } body { new_password }
	// binding request
	var bodyReq interface{}
	if err := handler.BindAndStructValidator(c, &bodyReq, false); err != nil {
		return
	}
	bodyReqMap := bodyReq.(map[string]interface{})
	newPassword := bodyReqMap["new_password"].(string)
	resetToken := c.Query("reset_token")

	// check is reset token exist
	resetpwData := models.ResetPassword{
		ResetToken: resetToken,
	}
	queryResetPassword := a.Db.Model(resetpwData).Where("reset_token = ?", resetToken)
	val, count := handler.QueryValidator(queryResetPassword, c, true)
	if !val && count == -1 {
		return
	}

	// if reset token not exist
	if count == 0 {
		handler.Error(c, http.StatusForbidden, utils.RESETPW_INVALID_TOKEN, "Invalid reset token")
		return
	}

	// if reset token exist
	queryResetPassword.First(&resetpwData)
	diff := utils.GetTimeDifference(resetpwData.UpdatedAt)

	// check is reset token expired
	if diff > 10 {
		handler.Error(c, http.StatusForbidden, utils.RESETPW_EXPIRED_TOKEN, "Reset token expired")
		return
	}

	// if reset token not expired
	// update password
	reqUser := models.User{
		UUID: resetpwData.UUID,
	}
	queryUser := a.Db.Model(reqUser)

	if val, _ = handler.QueryValidator(queryUser, c, false); !val {
		return
	}
	queryUser.First(&reqUser)
	encryptResult := utils.EncryptPassword(newPassword)
	if encryptResult == "" {
		handler.Error(c, http.StatusInternalServerError, utils.FAILED_ENCRYPT, "Encrypt Password Failed")
		return
	}
	reqUser.Password = encryptResult

	queryUpdate := a.Db.Model(reqUser).Updates(reqUser)
	if val, _ = handler.QueryValidator(queryUpdate, c, false); !val {
		return
	}

	// delete reset token
	queryDelete := a.Db.Delete(&models.ResetPassword{ID: resetpwData.ID})
	if val, _ = handler.QueryValidator(queryDelete, c, false); !val {
		return
	}

	handler.Success(c, http.StatusOK, "Reset Password Success", nil)
}

// utils
func sendResetTokenEmail(reqUser models.User, reqReset models.ResetPassword, token string) {
	handler.SendMail(handler.MailInfo{
		EmailTarget: reqUser.Email,
		Subject:     "Reset Password",
		Body: fmt.Sprintf(
			handler.ResetPassword,
			reqUser.Username,
			reqReset.Device,
			reqReset.IPAddress,
			utils.GetEnv("FE_URL")+"/resetpw/?reset_token="+token,
		),
	})
}
