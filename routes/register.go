package routes

import (
	"net/http"
	"strconv"

	"api.com/api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	user_id := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Coult not parse event id"})
		return
	}

	event, err := models.GetEventById(eventId)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	err = event.Register(user_id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Coult not register user for event"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Registered"})

}

func cancelRegisteration(context *gin.Context) {
	user_id := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Coult not parse event id"})
		return
	}

	var event models.Event
	event.ID = eventId

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	err = event.CancelRegister(user_id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Coult not cancel event id"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Cancelled"})

}
