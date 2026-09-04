package routes

import (
	"errors"
	"net/http"

	"api.com/api/models"
	"github.com/gin-gonic/gin"
)

func createUser(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = user.Save()

	if errors.Is(err, models.ErrUserExists) {
		context.JSON(http.StatusConflict, gin.H{"message": "username already taken"})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func getAllUsers(context *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch users"})
		return
	}
	context.JSON(http.StatusOK, users)
}
