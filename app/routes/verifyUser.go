package routes

import (
	"advanceauth/backend/app/controllers"
	"advanceauth/backend/app/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VerifyUserRoute(router *gin.Engine, db *gorm.DB) {
	verifyUserController := controller.VerifyUserController{
		Db: db,
	}
	verify := router.Group("/verify")

	verify.GET("/resend", middleware.CheckAuth(), verifyUserController.ResendVerify)
	verify.GET("/validate", verifyUserController.ValidateVerifyToken)
}
