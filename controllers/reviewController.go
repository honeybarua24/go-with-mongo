package controllers

import (
	"context"
	"log"
	"strconv"

	//"log"
	"go-with-mongo/database"
	helper "go-with-mongo/helper"
	models "go-with-mongo/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection = database.OpenCollection(database.Client, "review")

// Add  new review
func AddAReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "USER"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var review models.Reviews
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}
		//use the validator library to validate required fields
		if validationErr := validate.Struct(&review); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newReview := models.Reviews{
			Id:        primitive.NewObjectID(),
			Movie_id:  review.Movie_id,
			Review_id: review.Review_id,

			Review:      review.Review,
			Created_at:  review.Created_at,
			Updated_at:  review.Updated_at,
			Reviewer_id: review.Reviewer_id,
		}

		result, err := reviewCollection.InsertOne(ctx, newReview)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "success",
			"Data":    map[string]interface{}{"data": result}})
	}
}

// Filter reviews by movie_id
func ViewAMovieReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchreviews []models.Reviews
		movieId := c.Param("movie_id")
		i, erro := strconv.Atoi(movieId)
		if erro != nil {
			// Handle error
		}
		if movieId == "" {
			log.Println("No movie id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := reviewCollection.Find(ctx, bson.M{"movie_id": i})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchreviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchreviews)
	}
}

// Delete a review
func DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		reviewId := c.Param("review_id")
		defer cancel()
		i, erro := strconv.Atoi(reviewId)
		if erro != nil {
			// Handle error
		}

		result, err := reviewCollection.DeleteOne(ctx, bson.M{"review_id": i})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				gin.H{
					" Status":  http.StatusNotFound,
					" Message": "error",
					" Data":    map[string]interface{}{"data": "Review with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Your review was successfully deleted!"}},
		)
	}
}

// Allow a user view all their Reviews
func AllUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchreviews []models.Reviews
		reviewId := c.Param("reviewer_id")
		if reviewId == "" {
			log.Println("No reviewer id passed")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := reviewCollection.Find(ctx, bson.M{"reviewer_id": reviewId})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchreviews)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchreviews)
	}
}
