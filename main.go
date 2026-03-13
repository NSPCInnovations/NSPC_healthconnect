package main

import (
	"NSPC_healthconnect/database"
	"NSPC_healthconnect/doctor"
	"NSPC_healthconnect/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Connect Database
	database.ConnectDB()

	// Auto migrate all doctor related tables
	database.DB.AutoMigrate(
		&doctor.Doctor{},
		&doctor.DoctorVerification{},
		&doctor.DoctorAvailability{},
		&doctor.DoctorHospitalMap{},
		&doctor.DoctorDocument{},
		&doctor.DoctorRating{},
	)

	// Initialize Gin router
	r := gin.Default()

	// Register doctor routes
	router.DoctorRoutes(r)

	// Start server
	r.Run(":8080")
}
