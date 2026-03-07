package handlers

import (
	"net/http"

	"NSPC_healthconnect/database"
	"NSPC_healthconnect/doctor"

	"github.com/gin-gonic/gin"
)

func RegisterDoctor(c *gin.Context) {

	var doc doctor.Doctor

	// Bind JSON
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	// Validation check
	msg := ValidateDoctor(doc)

	if msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": msg,
		})
		return
	}

	// 🔴 Duplicate check
	var existing doctor.Doctor

	err := database.DB.
		Where(
			"doctor_reg_no = ? OR user_id = ? OR email = ? OR mobile = ?",
			doc.DoctorRegNo,
			doc.UserID,
			doc.Email,
			doc.Mobile,
		).
		First(&existing).Error

	if err == nil {

		c.JSON(http.StatusConflict, gin.H{
			"message": "Doctor already exists with same registration number, user id, email, or mobile",
		})
		return
	}

	// Default status
	doc.VerificationStatus = "PENDING"

	// Create doctor
	if err := database.DB.Create(&doc).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register doctor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor registered successfully",
		"data":    doc,
	})
}

func VerifyDoctor(c *gin.Context) {

	id := c.Param("id")

	var doc doctor.Doctor

	database.DB.First(&doc, id)

	doc.VerificationStatus = "APPROVED"

	database.DB.Save(&doc)

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor verified",
	})
}

func AddDoctorAvailability(c *gin.Context) {

	var availability doctor.DoctorAvailability

	c.ShouldBindJSON(&availability)

	database.DB.Create(&availability)

	c.JSON(http.StatusOK, gin.H{
		"message": "Availability added",
	})
}

func ListDoctors(c *gin.Context) {

	var doctors []doctor.Doctor

	database.DB.Where("verification_status = ?", "APPROVED").Find(&doctors)

	c.JSON(http.StatusOK, doctors)
}

func MapDoctorHospital(c *gin.Context) {

	var mapping doctor.DoctorHospitalMapping

	c.ShouldBindJSON(&mapping)

	database.DB.Create(&mapping)

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor mapped",
	})
}
