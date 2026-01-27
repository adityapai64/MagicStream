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
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

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