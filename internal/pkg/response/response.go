package response

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"net/http"
)

type ErrMessage struct {
	Message string `json:"message"` // 对外
	Inner   string `json:"inner"`
	Code    string `json:"code"`
}

func Error(c *gin.Context, httpCode int, message, inner string, code string) {
	c.JSON(httpCode, ErrMessage{Message: message, Inner: inner, Code: code})
}

func Success(c *gin.Context, obj interface{}) {
	if obj == nil {
		obj = gin.H{"message": "success"}
	}
	c.JSON(http.StatusOK, obj)
}

// 请求参数错误: 400
func BadRequest(c *gin.Context, message string, inner string) {
	Error(c, http.StatusBadRequest, message, inner, errors.CodeError)
}

// 认证出错: 401
func Unauthorized(c *gin.Context, message string, inner string) {
	Error(c, http.StatusUnauthorized, message, inner, errors.CodeError)
}

//系统内部错误: 500
func InternalErr(c *gin.Context, message string, inner string) {
	Error(c, http.StatusInternalServerError, message, inner, errors.CodeError)
}

// 请求参数错误: 返回 "请求参数错误"
func ParamsError(c *gin.Context, inner string) {
	BadRequest(c, "请求参数错误", inner)
}
