package routes

import (
	controllers "go-with-mongo/controllers"
	"go-with-mongo/middleware"

	"github.com/gin-gonic/gin"
)

func GenreRoutes(router gin.Engine) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/genres/creategenre", controllers.CreateGenre())
	router.GET("/genres/:genre_id", controllers.GetGenre())
	router.GET("/genres/name/:genre_name", controllers.SearchByName())
	router.GET("/genres", controllers.GetGenres())
	router.PUT("/genres/:genre_id", controllers.EditGenre())
	router.DELETE("/genres/:genre_id", controllers.DeleteGenre())
}
