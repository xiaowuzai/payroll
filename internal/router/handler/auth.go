package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	AccountName string `json:"accountName"` // 账户名称
	Password string `json:"password"`
}

// @Summary 登录
// @Description 登录验证
// @Tags 认证
// @Accept application/json
// @Param LoginRequest body LoginRequest true "accountName(required,len=6-16) password(required)"
// @Success 200 {object} ""
// @Router /login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	login := &LoginRequest{}
	err := c.ShouldBindJSON(login)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	token, refresh,err := uh.user.Login(ctx, login.AccountName,login.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK,gin.H{"token":token, "refreshToken": refresh})
}
