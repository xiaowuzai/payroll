package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/middleware"
	"github.com/xiaowuzai/payroll/internal/pkg/requestid"
	"github.com/xiaowuzai/payroll/internal/pkg/response"
	"net/http"
)

type LoginRequest struct {
	AccountName string `json:"accountName"` // 账户名称
	Password    string `json:"password"`
}

// @Summary 登录
// @Description 登录验证
// @Tags 认证
// @Accept application/json
// @Param LoginRequest body LoginRequest true "accountName(required,len=6-16) password(required)"
// @Success 200 {object} ""
// @Router /login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("Login function called")

	login := &LoginRequest{}
	err := c.ShouldBindJSON(login)
	if err != nil {
		log.Error("ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	token, refresh, err := uh.user.Login(ctx, login.AccountName, login.Password)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("Login function success")
	c.JSON(http.StatusOK, gin.H{"token": token, "refreshToken": refresh})
}

func (uh *UserHandler) WhoAmI(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("WhoAmI function called")

	authInfo, err := middleware.ParseJWT(c)
	if err != nil {
		log.Error("WhoAmI function ParseJWT error: ", err.Error())
		response.Unauthorized(c, err.Error(), err.Error())
		return
	}

	log.Info("WhoAmI function success")
	response.Success(c, authInfo)
}
