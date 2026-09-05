package routes

import (
	"api.com/api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.GET("/events", getEvents)
	authenticated.GET("/event/:id", getEventById)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/event/:id", updateEventById)
	authenticated.DELETE("/event/:id", deleteEventById)
	authenticated.GET("/users", getAllUsers)

	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegisteration)

	server.POST("/signup", createUser)
	server.POST("/login", login)
}
