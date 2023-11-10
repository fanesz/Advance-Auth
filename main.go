package main

import (
	"advanceauth/backend/app/config"
	"advanceauth/backend/app/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()

	router := gin.Default()
	router.Use(CORS())

	routes.AuthRoute(router, db)
	routes.UserRoute(router, db)
	routes.ResetPasswordRoute(router, db)
	routes.VerifyUserRoute(router, db)

	router.Run("localhost:8080")
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	}
}
