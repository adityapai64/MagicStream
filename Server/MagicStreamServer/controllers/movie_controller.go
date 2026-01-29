package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/adityapai64/MagicStream/Server/MagicStreamServer/database"
	"github.com/adityapai64/MagicStream/Server/MagicStreamServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"github.com/go-playground/validator/v10"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var validate = validator.New()

/*
** Function names which begin with capital letters are public. If the first letter
** of the function name is not capitalised, the function has a package protected
** visibility.
 */
func GetMovies() gin.HandlerFunc {
	return  func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies [] models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if (err != nil) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode movies"})
		}
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		movieID := c.Param("imdb_id")

		if (movieID == "") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return 
		}
		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if (err != nil) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie

		/*
			What's happening in the code below:

			- the parameter c holds the context of the http request (line 68)
			- the method c.shouldBindJSON(&movie) tries to map the JSON object sent 
			through the request to the movie variable, declared on line 72
			- If the mapping succeeds, the method returns null (or nil in go);
			if the mapping throws an error, the error object will be defined
			- Line 86 simply returns an http 400 error with an error message.
			- The & character (line 86) is the address of operator in go: it gets
			the memory address of the variable movie, i.e. it passes a pointer to 
			the variable movie and not its value
		*/

		if err:= c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		/*
			Validation:

			- Validation is done using the go-playground/validator library
			- A struct's attrubutes are bound to validation constraints using back ticks (Check
			the structs in movie_model.go) 
		*/
		if err:= validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		
		/*
		*The code below finally inserts an object (a movie into the database)
		* - movieCollection is a MongoDB collection object that contains all 
		* 	movie data.
		* - The insertOne function inserts one document into the collection
		* - The first parameter of insertOne, ctx, is the context used for 
		*	timeout, cancellation or passing metadata
		* - Movie is the Go Struct that will be converted into a BSON
		*	document and inserted into the collection
		*/

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add movie"})
			return 
		}

		c.JSON(http.StatusCreated, result)
	}
}