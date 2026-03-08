package hospital

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	c.JSON(200, response)
}

func ErrorResponse(c *gin.Context, message string) {
	response := APIResponse{
		Status:  "error",
		Message: message,
	}

	c.JSON(400, response)
}
