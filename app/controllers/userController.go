package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"advanceauth/backend/app/handler"
	"advanceauth/backend/app/models"
	"advanceauth/backend/app/utils"
	"net/http"
)

type UserController struct {
	Db *gorm.DB
}

func (a *UserController) Register(c *gin.Context) {
	// body { username, email, password }
	// binding and validating request
	req := models.User{}
	if err := handler.BindAndStructValidator(c, &req, true); err != nil {
		return
	}

	// encrypt password
	encryptResult := utils.EncryptPassword(req.Password)
	if encryptResult == "" {
		handler.Error(c, http.StatusInternalServerError, utils.FAILED_ENCRYPT, encryptResult)
		return
	}
	req.Password = encryptResult

	// create user
	user := models.User{
		UUID:     uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	queryCreate := a.Db.Create(&user)
	val, _ := handler.QueryValidator(queryCreate, c, false)
	if !val {
		return
	}

	handler.Success(c, http.StatusOK, "Create Success", nil)
}

func (a *UserController) UpdateUsername(c *gin.Context) {
	// body { username, device, ip_address }
	// binding and validating request
	var bodyReq interface{}
	if err := handler.BindAndStructValidator(c, &bodyReq, false); err != nil {
		return
	}
	bodyReqMap := bodyReq.(map[string]interface{})
	reqUser := models.User{
		Username: bodyReqMap["username"].(string),
	}
	reqInfo := models.LoggedUser{
		IPAddress: bodyReqMap["ip_address"].(string),
		Device:    bodyReqMap["device"].(string),
	}

	// checking is user logged in
	authController := AuthController{
		Db: a.Db,
	}
	if check := authController.LoginChecker(c, reqInfo); !check {
		return
	}

	// update username
	query := a.Db.Model(&models.User{UUID: c.GetString("UUID")}).Updates(reqUser)
	val, _ := handler.QueryValidator(query, c, false)
	if !val {
		return
	}

	handler.Success(c, http.StatusOK, "Update Username Success", gin.H{"username": reqUser.Username})
}
