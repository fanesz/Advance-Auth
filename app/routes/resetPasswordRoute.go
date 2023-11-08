package routes

import (
	"advanceauth/backend/app/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ResetPasswordRoute(router *gin.Engine, db *gorm.DB) {
	ResetPwController := controller.ResetPwController{
		Db: db,
	}
	resetpw := router.Group("/resetpw")

	resetpw.POST("/request", ResetPwController.ResetRequest)
	resetpw.POST("/validate", ResetPwController.ValidateResetRequest)
	resetpw.POST("/reset", ResetPwController.ResetPassword)
}
