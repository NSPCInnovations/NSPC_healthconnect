package hospital

import (
	"NSPC_HEALTHCONNECT/database"

	"github.com/gin-gonic/gin"
)

func RegisterHospital(c *gin.Context) {

	var hospital Hospital

	// JSON request read
	if err := c.ShouldBindJSON(&hospital); err != nil {
		ErrorResponse(c, "Invalid request data")
		return
	}

	// default status
	hospital.Status = "Pending"

	// save to database
	result := database.DB.Create(&hospital)

	if result.Error != nil {
		ErrorResponse(c, "Failed to register hospital")
		return
	}

	SuccessResponse(c, "Hospital registered successfully", hospital)
}
