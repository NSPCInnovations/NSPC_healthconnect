package router

import (
	"NSPC_healthconnect/handlers"

	"github.com/gin-gonic/gin"
)

func DoctorRoutes(r *gin.Engine) {

	group := r.Group("/api/doctors")

	// Doctor Registration
	group.POST("/register", handlers.RegisterDoctor)

	// Verify Doctor
	group.PUT("/verify/:id", handlers.VerifyDoctor)

	// Add Availability
	group.POST("/availability", handlers.AddDoctorAvailability)

	// Map Doctor Hospital
	group.POST("/map-hospital", handlers.MapDoctorHospital)

	// Upload Document
	group.POST("/document", handlers.UploadDoctorDocument)

	// Add Rating
	group.POST("/rating", handlers.AddDoctorRating)

	// List Doctors
	group.GET("/list", handlers.ListDoctors)
}
