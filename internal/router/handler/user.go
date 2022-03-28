package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/service"
	"net/http"
)

type UserHandler struct {
	user *service.UserService
}
func NewUserHandler(userService *service.UserService) *UserHandler{
	return &UserHandler{
		user: userService,
	}
}

type User struct {
	Id string `json:"id"`
	Username string `json:"username" binding:"required"`
	AccountName string `json:"accountName" binding:"required"`
	Email string `json:"email"`
	RoleId string `json:"roleId" binding:"required"`
	Password string `json:"password" binding:"required"`
	Status int32 `json:"status"`
	Created uint64 `json:"created"`
}

// @Summary 添加用户
// @Description 添加用户时指定角色
// @Tags 用户管理
// @Accept application/json
// @Param User body User true ""
// @Success 200 {object} response.Result ""
// @Router /v1/auth/user [post]
func (uh *UserHandler)AddUser(c *gin.Context) {
	user :=  &User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	err = uh.user.AddUser(ctx, &service.User{
		Username: user.Username,
		AccountName: user.AccountName,
		Email: user.Email,
		RoleId: user.RoleId,
		Password: user.Password,
		Status: user.Status,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (uh *UserHandler)AddAdmin(c *gin.Context) {
	ctx := c.Request.Context()
	err := uh.user.AddUser(ctx, &service.User{
		Username: "管理员",
		AccountName: "admin",
		Email: "",
		RoleId: "4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",
		Password: "Admin@123",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
