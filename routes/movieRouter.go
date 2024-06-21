package routes

import (
	controllers "go-with-mongo/controllers"
	"go-with-mongo/middleware"

	"github.com/gin-gonic/gin"
)

func MovieRoutes(router gin.Engine) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/movies/createmovie", controllers.CreateMovie())
	router.GET("/movies/:movie_id", controllers.GetMovie())
	router.GET("/movies", controllers.GetMovies())
	router.PUT("/movies/:movie_id", controllers.UpdateMovie())
	router.GET("/movies/search/:movieName", controllers.SearchMovieByQuery())
	router.GET("/movies/filter/:genreID", controllers.SearchMovieByGenre())
	router.DELETE("/movies/:movie_id", controllers.DeleteMovie())
}
