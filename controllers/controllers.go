package controllers

import (
	"MongoService/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if models.InsertUser(user) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User inserted successfully!"})
}

func InsertUsers(c *gin.Context) {
	var users []models.User
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.InsertUsers(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert users"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Users inserted successfully!"})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully!"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := models.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully!"})
}

func DeleteAll(c *gin.Context) {
	err := models.DeleteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All users deleted successfully!"})
}

func ListAllUsers(c *gin.Context) {
	users, err := models.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, users)

}

func FindUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := models.FindUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}
