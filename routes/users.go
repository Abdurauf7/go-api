package routes

import (
	"errors"
	"net/http"

	"api.com/api/models"
	"api.com/api/utils"
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

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = user.ValidateCredentials()

	if errors.Is(err, models.ErrInvalidCredentials) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not authenticate user"})
		return
	}
	token, err := utils.GenerateToken(user.Username, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login in success", "token": token})
}
