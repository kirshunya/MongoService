package main

import (
	"MongoService/controllers"
	"MongoService/models"
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

	models.ConnectMongo()

	userGroup := r.Group("/users")
	{
		userGroup.POST("/create", controllers.InsertUser)
		userGroup.POST("/create_all", controllers.InsertUsers)
		userGroup.PUT("/update/:id", controllers.UpdateUser)
		userGroup.DELETE("/delete/:id", controllers.DeleteUser)
		userGroup.GET("/:id", controllers.FindUserById)
		userGroup.GET("/all", controllers.ListAllUsers)
		userGroup.DELETE("/delete/all", controllers.DeleteAll)
	}

	r.Run(":8081")

}
