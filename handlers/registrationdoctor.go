package handlers

import (
	"net/http"

	"NSPC_healthconnect/database"
	"NSPC_healthconnect/doctor"

	"github.com/gin-gonic/gin"
)

func RegisterDoctor(c *gin.Context) {

	var doc doctor.Doctor

	if err := c.ShouldBindJSON(&doc); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON",
		})

		return
	}

	msg := ValidateDoctor(doc)

	if msg != "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": msg,
		})

		return
	}

	doc.VerificationStatus = "PENDING"

	database.DB.Create(&doc)

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor registered",
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
