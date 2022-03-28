package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}

func Success(c *gin.Context, obj interface{}) {
	if obj == nil {
		obj = gin.H{"message": "success"}
	}
	c.JSON(http.StatusOK, obj)
}