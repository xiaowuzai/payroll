package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuccessMessage struct{
	Message string `json:"message"`
}

var success string = "success"


func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}

func Success(c *gin.Context, obj interface{}) {
	if obj == nil {
		obj = SuccessMessage{
			Message: success,
		}
	}
	c.JSON(http.StatusOK, obj)
}

type BadMessage struct {
	Message string `json:"message"`
	InnerMessage string `json:""`
}

// 请求参数错误: 400
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// 认证出错: 401
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized,message)
}

//系统内部错误: 500
func InternalErr(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
