package handlers

import (
	"net/http"
	"time"

	"NSPC_healthconnect/database"
	"NSPC_healthconnect/doctor"

	"github.com/gin-gonic/gin"
)

func RegisterDoctor(c *gin.Context) {

	var doc doctor.Doctor

	if err := c.ShouldBindJSON(&doc); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
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

	var existing doctor.Doctor

	err := database.DB.
		Where("doctor_reg_no = ? OR user_id = ? OR email = ? OR mobile = ?",
			doc.DoctorRegNo,
			doc.UserID,
			doc.Email,
			doc.Mobile).
		First(&existing).Error

	if err == nil {

		c.JSON(http.StatusConflict, gin.H{
			"message": "Doctor already exists",
		})
		return
	}

	doc.VerificationStatus = "PENDING"
	doc.IsActive = false

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

	if err := database.DB.First(&doc, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Doctor not found",
		})
		return
	}

	doc.VerificationStatus = "APPROVED"
	doc.IsActive = true

	database.DB.Save(&doc)

	verification := doctor.DoctorVerification{
		DoctorID:           doc.ID,
		VerificationStatus: "APPROVED",
		VerifiedBy:         "ADMIN",
		VerifiedAt:         time.Now(),
	}

	database.DB.Create(&verification)

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor verified successfully",
	})
}

func AddDoctorAvailability(c *gin.Context) {

	var availability doctor.DoctorAvailability

	if err := c.ShouldBindJSON(&availability); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := database.DB.Create(&availability).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add availability",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor availability added",
	})
}
func MapDoctorHospital(c *gin.Context) {

	var mapping doctor.DoctorHospitalMap

	if err := c.ShouldBindJSON(&mapping); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := database.DB.Create(&mapping).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to map doctor to hospital",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor mapped to hospital",
	})
}

func UploadDoctorDocument(c *gin.Context) {

	var doc doctor.DoctorDocument

	if err := c.ShouldBindJSON(&doc); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	var doctorExists doctor.Doctor

	if err := database.DB.First(&doctorExists, doc.DoctorID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Doctor not found",
		})
		return
	}

	if err := database.DB.Create(&doc).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload document",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor document uploaded successfully",
	})
}

func AddDoctorRating(c *gin.Context) {

	var rating doctor.DoctorRating

	if err := c.ShouldBindJSON(&rating); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var doc doctor.Doctor

	if err := database.DB.First(&doc, rating.DoctorID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Doctor not found",
		})
		return
	}

	if err := database.DB.Create(&rating).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add rating",
		})
		return
	}

	var avg float64

	database.DB.Model(&doctor.DoctorRating{}).
		Where("doctor_id = ?", rating.DoctorID).
		Select("AVG(rating)").Scan(&avg)

	database.DB.Model(&doctor.Doctor{}).
		Where("id = ?", rating.DoctorID).
		Update("average_rating", avg)

	c.JSON(http.StatusOK, gin.H{
		"message": "Rating added successfully",
	})
}
func ListDoctors(c *gin.Context) {

	var doctors []doctor.Doctor

	if err := database.DB.
		Where("verification_status = ?", "ACTIVE").
		Find(&doctors).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch doctors",
		})
		return
	}

	c.JSON(http.StatusOK, doctors)
}