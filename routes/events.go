package routes

import (
	"errors"
	"net/http"
	"strconv"

	"api.com/api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch events"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user_id := context.GetInt64("user_id")
	event.UserId = user_id

	if err := event.Save(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "event created", "event": event})
}

func getEventById(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse event id"})
		return
	}

	event, err := models.GetEventById(eventId)

	if errors.Is(err, models.ErrEventNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch event"})
		return
	}

	context.JSON(http.StatusOK, event)

}

func updateEventById(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse event id"})
		return
	}

	user_id := context.GetInt64("user_id")
	existingEvent, err := models.GetEventById(eventId)

	if errors.Is(err, models.ErrEventNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.UserId != user_id {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized to update event "})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch event"})
		return
	}

	var updatedEvent models.Event

	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	updatedEvent.ID = eventId
	updatedEvent.UserId = existingEvent.UserId

	if err := updatedEvent.Update(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "event updated", "event": updatedEvent})
}

func deleteEventById(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse event id"})
		return
	}

	event, err := models.GetEventById(eventId)

	user_id := context.GetInt64("user_id")

	if errors.Is(err, models.ErrEventNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.UserId != user_id {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized to update event "})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch event"})
		return
	}

	if err := event.Delete(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}
