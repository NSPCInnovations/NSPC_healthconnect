package router

import (
	"NSPC_healthconnect/handlers"

	"github.com/gin-gonic/gin"
)

func DoctorRoutes(router *gin.Engine) {

	group := router.Group("/api/doctors")

	group.POST("/register", handlers.RegisterDoctor)
	group.PUT("/verify/:id", handlers.VerifyDoctor)
	group.POST("/availability", handlers.AddDoctorAvailability)
	group.GET("/list", handlers.ListDoctors)
	group.POST("/map-hospital", handlers.MapDoctorHospital)
}
