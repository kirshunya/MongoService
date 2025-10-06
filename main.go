package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	user := r.Group("/user")
	{
		user.POST("/create")
	}

	r.Run()

}
