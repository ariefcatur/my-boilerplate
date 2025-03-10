package helpers

import "github.com/gin-gonic/gin"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func APIResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	resp := Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}

	c.JSON(statusCode, resp)
}
func ErrorResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	resp := Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}

	c.JSON(statusCode, resp)
}
