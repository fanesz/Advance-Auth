package routes

import (
	"advanceauth/backend/app/controllers"
	"advanceauth/backend/app/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoute(router *gin.Engine, db *gorm.DB) {
	authController := controller.AuthController{
		Db: db,
	}
	auth := router.Group("/auth")

	auth.POST("/login", authController.Login)
	auth.GET("/isLogin", middleware.CheckAuth(), authController.IsLogin)
	auth.GET("/logout", middleware.CheckAuth(), authController.Logout)
}
