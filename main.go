package main

import (
	"api.com/api/db"
	"api.com/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	server := gin.Default()
	routes.RegisterRoutes(server)

	server.Run(":8080")
}
