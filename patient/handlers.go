package patient

import (
	"net/http"

	"NSPC_HEALTHCONNECT/database"

	"github.com/gin-gonic/gin"
)

func RegisterPatient(c *gin.Context) {

	var patient Patient

	// 1. Read JSON request
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILED",
			"message": "Invalid request body",
		})
		return
	}
	if patient.FirstName == "" || patient.LastName == "" ||
		patient.MobileNumber == "" || patient.Age <= 0 || patient.Gender == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILED",
			"message": "Missing or invalid required fields",
		})
		return
	}
	var existingPatient Patient

	err := database.DB.
		Where("mobile_number = ?", patient.MobileNumber).
		First(&existingPatient).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "FAILED",
			"message": "Mobile number already exists",
		})
		return
	}
	if err := database.DB.Create(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": "Failed to register patient",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":     "SUCCESS",
		"message":    "Patient registered successfully",
		"patient_id": patient.PatientID,
	})
}
