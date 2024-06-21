package routes

import (
	"go-with-mongo/controllers"
	"go-with-mongo/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router gin.Engine) {
	router.Use(middleware.AuthenticateUser())
	router.GET("/users/:user_id", controllers.GetUser())
	router.GET("/users", controllers.GetUsers())
	router.POST("/users", controllers.CreateUser())

}
