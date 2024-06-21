package routes

import (
	controllers "go-with-mongo/controllers"
	"go-with-mongo/middleware"

	"github.com/gin-gonic/gin"
)

func ReviewRoutes(router gin.Engine) {
	router.Use(middleware.AuthenticateUser())
	router.POST("reviews/addreview", controllers.AddAReview())
	router.GET("/reviews/filter/:movie_id", controllers.ViewAMovieReviews())
	router.DELETE("reviews/:review_id", controllers.DeleteReview())
	router.GET("reviews/user_reviews/:reviewer_id", controllers.AllUserReviews())
}
