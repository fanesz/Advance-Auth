package routes

import (
	"advanceauth/backend/app/controllers"
	"advanceauth/backend/app/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoute(router *gin.Engine, db *gorm.DB) {
	userController := controller.UserController{
		Db: db,
	}
	user := router.Group("/user")

	user.POST("/register", userController.Register)
	user.PATCH("/update/username", middleware.CheckAuth(), userController.UpdateUsername)
}
