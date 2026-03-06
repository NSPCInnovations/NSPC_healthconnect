package main

import (
	"fmt"
	"log"

	"NSPC_HEALTHCONNECT/database"
	"NSPC_HEALTHCONNECT/patient"

	"github.com/gin-gonic/gin"
)

func main() {

	// 1️⃣ DATABASE CONNECTION CHECK
	err := database.ConnectDatabase()
	if err != nil {
		log.Println(" Database connection failed")
		log.Fatal(err)
	} else {
		fmt.Println(" Database connected successfully")
	}

	// 2️⃣ AUTO MIGRATION CHECK
	if err := database.DB.AutoMigrate(&patient.Patient{}); err != nil {
		log.Println(" Auto migration failed")
		log.Fatal(err)
	} else {
		fmt.Println(" Patient table migrated successfully")
	}

	// 3️⃣ SERVER START
	router := gin.Default()
	patient.RegisterPatientRoutes(router)

	fmt.Println(" Server running on http://localhost:8080")
	router.Run(":8080")
}
