package main

import (
	"os"

	"nspc_healthcare/config"
	"nspc_healthcare/doctor"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// load env
	godotenv.Load()

	// db connect
	config.ConnectDB()

	// migrate
	config.DB.AutoMigrate(&doctor.Doctor{})

	r := gin.Default()

	// 🔥 TEST ROUTE
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// doctor routes
	doctor.RegisterDoctorRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
