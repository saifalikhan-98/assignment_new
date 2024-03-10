// cmd/main.go

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"yourproject/pkg/user"
)

func main() {
	// Initialize the user manager (you can choose between InMemoryUserManager or SQLiteUserManager)
	// Example with InMemoryUserManager:
	userManager := user.NewConcurrentUserManager()

	// Example with SQLiteUserManager (uncomment this line and comment the previous line):
	// userManager, err := user.NewSQLiteUserManager("sqlite.db")
	// if err != nil {
	// 	fmt.Printf("Error initializing SQLiteUserManager: %v\n", err)
	// 	return
	// }

	router := gin.Default()

	// Define RESTful endpoints
	router.POST("/users", CreateUserHandler(userManager))
	router.GET("/users/:id", ReadUserHandler(userManager))
	router.PUT("/users/:id", UpdateUserHandler(userManager))
	router.DELETE("/users/:id", DeleteUserHandler(userManager))

	// Start the server
	router.Run(":8080")
}

func CreateUserHandler(manager user.UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser user.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := manager.Create(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": userID})
	}
}

func ReadUserHandler(manager user.UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		foundUser, err := manager.Read(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func UpdateUserHandler(manager user.UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var updatedUser user.User
		if err := c.ShouldBindJSON(&updatedUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := manager.Update(userID, updatedUser); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

func DeleteUserHandler(manager user.UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if err := manager.Delete(userID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
