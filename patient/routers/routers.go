package patient_registration

import "github.com/gin-gonic/gin"

func RegisterPatientRoutes(router *gin.Engine) {

	patientGroup := router.Group("/patients")
	{
		patientGroup.POST("/register", RegisterPatient)
	}

}
