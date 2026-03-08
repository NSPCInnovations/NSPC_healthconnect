package main

import (
	"NSPC_HEALTHCONNECT/database"
	"NSPC_HEALTHCONNECT/hospital"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()
	database.DB.AutoMigrate(&hospital.Hospital{})

	router := gin.Default()

	hospital.HospitalRoutes(router)

	router.Run(":8080")
}
