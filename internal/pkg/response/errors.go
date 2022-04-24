package response

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"net/http"
)

// 如果是 MyError 类型，则按照 MyError类型处理。否则返回 500 以及错误信息
func WithError(c *gin.Context, err error) {
	var myErr *errors.MyError
	if errors.As(err, &myErr) {
		e := err.(*errors.MyError)
		Error(c, e.HTTPCode(), e.Message(), e.Inner(), e.Code())
		return
	}

	Error(c, http.StatusInternalServerError, err.Error(), err.Error(), errors.CodeError)
	return
}
