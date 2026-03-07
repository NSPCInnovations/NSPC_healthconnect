package main

import (
	"NSPC_healthconnect/database"
	"NSPC_healthconnect/doctor"
	"NSPC_healthconnect/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	database.ConnectDB()

	database.DB.AutoMigrate(
		&doctor.Doctor{},
		&doctor.DoctorVerification{},
		&doctor.DoctorAvailability{},
		&doctor.DoctorHospitalMapping{},
	)

	r := gin.Default()

	router.DoctorRoutes(r)

	r.Run(":8080")
}
