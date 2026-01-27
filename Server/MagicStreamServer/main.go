package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	controller "github.com/adityapai64/MagicStream/Server/MagicStreamServer/controllers"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello from MagicStream!!")
	})

	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetMovie())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}
