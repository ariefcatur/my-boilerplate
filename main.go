package main

import (
	"github.com/ariefcatur/my-boilerplate/config"
	"github.com/ariefcatur/my-boilerplate/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Auth routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Run the server
	r.Run(":8080")

}
