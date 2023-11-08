package controller

import (
	"advanceauth/backend/app/handler"
	"advanceauth/backend/app/models"
	"advanceauth/backend/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

type AuthController struct {
	Db *gorm.DB
}

func (a *AuthController) Login(c *gin.Context) {
	// body { email, password, ip_address, device }
	// binding request
	var bodyReq interface{}
	if err := handler.BindAndStructValidator(c, &bodyReq, false); err != nil {
		return
	}
	bodyReqMap := bodyReq.(map[string]interface{})
	reqAuth := models.User{
		Email:    bodyReqMap["email"].(string),
		Password: bodyReqMap["password"].(string),
	}
	reqInfo := models.LoggedUser{
		IPAddress: bodyReqMap["ip_address"].(string),
		Device:    bodyReqMap["device"].(string),
	}
	validateLogin := models.User{
		Email: reqAuth.Email,
	}

	// get user by email
	queryUser := a.Db.Model(validateLogin).Where("email = ?", validateLogin.Email)
	val, count := handler.QueryValidator(queryUser, c, true)
	if !val && count == -1 {
		return
	}

	// if user not exist
	if count == 0 {
		handler.Error(c, http.StatusUnauthorized, utils.AUTH_EMAIL_NOT_FOUND, "Unauthorized")
		return
	}

	// if user exist
	queryUser.First(&validateLogin)

	// validate password
	comparingPassword := utils.DecryptPassword(validateLogin.Password, reqAuth.Password)
	if !comparingPassword {
		handler.Error(c, http.StatusUnauthorized, utils.AUTH_WRONG_PASSWORD, "Unauthorized")
		return
	}

	// checking is user already logged in
	storedLoggedUser := models.LoggedUser{
		UUID: validateLogin.UUID,
	}
	query := a.Db.Model(storedLoggedUser)
	val, count = handler.QueryValidator(query, c, true)
	if !val && count == -1 {
		return
	}

	// generate token
	validateToken := utils.PayloadToken{
		UUID: validateLogin.UUID,
	}
	newToken, err := utils.GenerateToken(&validateToken)
	if err != nil {
		handler.Error(c, http.StatusInternalServerError, utils.FAILED_GENERATE_JWT, err.Error())
		return
	}
	loggedUser := models.LoggedUser{
		LoginToken: newToken,
		IPAddress:  reqInfo.IPAddress,
		Device:     reqInfo.Device,
		UUID:       validateLogin.UUID,
	}

	// when user already logged in
	if count != 0 {
		query.Where("ip_address = ? AND device = ?", reqInfo.IPAddress, reqInfo.Device)
		val, count = handler.QueryValidator(query, c, true)
		if !val && count == -1 {
			return
		}
		// return logged user token if already logged in
		if count != 0 {
			query.First(&storedLoggedUser)
			handler.Success(c, http.StatusOK, "Login Success", gin.H{"token": storedLoggedUser.LoginToken})
			return
		}

		// when login from new device/browser/ipaddress
		if count == 0 {
			handler.SendMail(handler.MailInfo{
				EmailTarget: validateLogin.Email,
				Subject:     "New Device Login Login",
				Body: fmt.Sprintf(
					handler.NewDeviceLogin,
					loggedUser.Device,
					loggedUser.IPAddress,
				),
			})
			count = 0
		}
	}

	// when user not logged in
	if count == 0 {
		// create logged user record
		queryCreate := a.Db.Create(&loggedUser)
		val, _ := handler.QueryValidator(queryCreate, c, false)
		if !val {
			return
		}
		handler.Success(c, http.StatusOK, "Login Success", gin.H{"token": newToken})
	}
}

func (a *AuthController) IsLogin(c *gin.Context) {
	// body { ip_address, device }
	req := models.LoggedUser{}
	if err := handler.BindAndStructValidator(c, &req, false); err != nil {
		return
	}

	// checking is user logged in
	if check := a.LoginChecker(c, req); !check {
		return
	}

	handler.Success(c, http.StatusOK, "Is Login", gin.H{"UUID": c.GetString("UUID")})
}

func (a *AuthController) Logout(c *gin.Context) {
	// get context from middleware
	UUID := c.GetString("UUID")
	loginToken := c.GetString("loginToken")

	// checking is user logged in
	loggedUser := models.LoggedUser{}
	query := a.Db.Model(loggedUser).Where("login_token = ? AND UUID = ?", loginToken, UUID)
	val, count := handler.QueryValidator(query, c, true)
	if !val && count == -1 {
		return
	}

	// when user already logout
	if count == 0 {
		handler.Success(c, http.StatusOK, "User Already Logout", nil)
		return
	}

	// when user logged in
	if count != 0 {
		query.Delete(loggedUser)
	}

	handler.Success(c, http.StatusOK, "Logout Success", nil)
}

// utils
func (a *AuthController) LoginChecker(c *gin.Context, req models.LoggedUser) bool {
	// get context from middleware
	UUID := c.GetString("UUID")
	loginToken := c.GetString("loginToken")
	isTokenExpired := c.GetBool("isTokenExpired")

	// check is token expired
	if isTokenExpired {
		queryDelete := a.Db.Delete(&models.LoggedUser{}, "login_token = ?", loginToken)
		val, _ := handler.QueryValidator(queryDelete, c, false)
		if !val {
			return false
		}
		handler.Error(c, http.StatusUnauthorized, utils.AUTH_TOKEN_EXPIRED, "Auth Token Expired")
		return false
	}

	loggedUser := models.LoggedUser{}
	query := a.Db.Model(loggedUser).Where("login_token = ? AND UUID = ?", loginToken, UUID)
	val, count := handler.QueryValidator(query, c, true)
	if !val && count == -1 {
		return false
	}

	// when user not logged in
	if count == 0 {
		handler.Error(c, http.StatusUnauthorized, utils.AUTH_USER_NOT_LOGGEDIN, "Unauthorized")
		return false
	}

	// when user already logged in
	query.First(&loggedUser)

	// checking is user token got hacked
	if loggedUser.IPAddress != req.IPAddress || loggedUser.Device != req.Device {
		user := models.User{}
		queryUser := a.Db.Model(user).Where("uuid = ?", UUID)
		val, count := handler.QueryValidator(queryUser, c, true)
		if !val && count == -1 {
			return false
		}
		if count != 0 {
			queryUser.First(&user)
			handler.SendMail(handler.MailInfo{
				EmailTarget: user.Email,
				Subject:     "Unauthorized Login",
				Body: fmt.Sprintf(
					handler.LoginTokenGotHacked,
					req.Device,
					req.IPAddress,
				),
			})
			query.Delete(loggedUser)
			handler.Error(c, http.StatusUnauthorized, utils.AUTH_DIFFERENT_IP_OR_DEVICE, "Unauthorized")
			return false
		}
	}

	return true
}
