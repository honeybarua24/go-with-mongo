package controllers

import (
	"context"
	"go-with-mongo/database"
	helper "go-with-mongo/helper"
	models "go-with-mongo/model"
	"log"
	"net/http"
	"strconv"
	"time"

	//"strconv"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection(database.Client, "movie")

// To create one movie
func CreateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var movie models.Movie
		defer cancel()

		if err := c.BindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		//Check to see if name exists
		regexMatch := bson.M{"$regex": primitive.Regex{Pattern: *movie.Name, Options: "i"}}
		count, err := movieCollection.CountDocuments(ctx, bson.M{"name": regexMatch})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the movie name"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this movie name already exists", "count": count})
			return
		}

		if validationError := validate.Struct(&movie); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}
		newMovie := models.Movie{
			Id:       primitive.NewObjectID(),
			Name:     movie.Name,
			Topic:    movie.Topic,
			Genre_id: movie.Genre_id,

			Movie_URL:  movie.Movie_URL,
			Created_at: movie.Created_at,
			Updated_at: movie.Updated_at,
			Movie_id:   movie.Movie_id,
		}
		err = movieCollection.FindOne(ctx, bson.M{"movie_id": movie.Movie_id}).Err()
		defer cancel()
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this movie_id already exists"})
			return
		} else if err != mongo.ErrNoDocuments {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the movie_id"})
			return
		}
		result, err := movieCollection.InsertOne(ctx, newMovie)

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

// To get just one movie
func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		movieId := c.Param("movie_id")
		var movie models.Movie
		defer cancel()

		i, erro := strconv.Atoi(movieId)
		if erro != nil {
			// Handle error
		}
		err := movieCollection.FindOne(ctx, bson.M{"movie_id": i}).Decode(&movie)

		//objId, _ := primitive.ObjectIDFromHex(movieId)

		//err := movieCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    map[string]interface{}{"data": movie}})
	}
}

// To fetch all movies
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "movie_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := movieCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching movies "})
		}
		var allmovies []bson.M
		if err = result.All(ctx, &allmovies); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allmovies[0])
	}
}

// Update a movie
func UpdateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		movieId := c.Param("movie_id")
		var movie models.Movie
		defer cancel()
		//objId, _ := primitive.ObjectIDFromHex(movieId)
		i, erro := strconv.Atoi(movieId)
		if erro != nil {
			// Handle error
		}
		if err := c.BindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationError := validate.Struct(&movie); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		update := bson.M{
			"name":      movie.Name,
			"topic":     movie.Topic,
			"genre_id":  movie.Genre_id,
			"movie_url": movie.Movie_URL}
		filterByID := bson.M{"movie_id": i}
		result, err := movieCollection.UpdateOne(ctx, filterByID, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"Status":  http.StatusNotFound,
				"Message": "movie not found",
				"Data":    nil})
			return
		}
		var updatedMovie models.Movie
		if result.MatchedCount == 1 {
			err := movieCollection.FindOne(ctx, filterByID).Decode(&updatedMovie)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Status":  http.StatusInternalServerError,
					"Message": "error",
					"Data":    map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "movie updated successfully!",
			"Data":    updatedMovie})
	}
}

func SearchMovieByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchmovies []models.Movie
		movieName := c.Param("movieName")

		if movieName == "" {
			log.Println("movie name is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search parameter" + movieName})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := movieCollection.Find(ctx, bson.M{"name": bson.M{"$regex": movieName}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong")
			return
		}
		err = searchquerydb.All(ctx, &searchmovies)
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
		c.IndentedJSON(200, searchmovies)

	}
}

// Filter movie by genre
func SearchMovieByGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchbygenre []models.Movie
		queryParam, err := strconv.Atoi(c.Param("genreID"))
		//c.IndentedJSON(200, gin.H{"genre_id": queryParam, "results": searchbygenre})

		if err != nil {
			log.Println("Invalid genre_id")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid genre_id"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchdb, err := movieCollection.Find(ctx, bson.M{"genre_id": queryParam})

		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchdb.All(ctx, &searchbygenre)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchdb.Close(ctx)
		if err := searchdb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchbygenre)

	}
}

func DeleteMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		movieId := c.Param("movie_id")
		defer cancel()
		i, erro := strconv.Atoi(movieId)
		if erro != nil {
			// Handle error
		}
		result, err := movieCollection.DeleteOne(ctx, bson.M{"movie_id": i})
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
					" Data":    map[string]interface{}{"data": "Movie with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Movie successfully deleted!"}},
		)
	}
}
