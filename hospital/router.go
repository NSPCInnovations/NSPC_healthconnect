package hospital

import "github.com/gin-gonic/gin"

func HospitalRoutes(router *gin.Engine) {

	router.POST("/hospital/register", RegisterHospital)

}
