package doctor

import (
	"nspc_healthcare/config"

	"github.com/gin-gonic/gin"
)

func RegisterDoctor(c *gin.Context) {

	var doctor Doctor

	//  SAFE JSON READ (works for application/json AND text/plain)
	err := c.ShouldBindJSON(&doctor)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  0,
			"message": "Invalid JSON data",
			"error":   err.Error(),
		})
		return
	}

	doctor.VerificationStatus = "PENDING"
	doctor.IsActive = false

	if err := config.DB.Create(&doctor).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  0,
			"message": "Doctor registration failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  1,
		"message": "Doctor registered successfully",
		"data":    doctor,
	})
}
