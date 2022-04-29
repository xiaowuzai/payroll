package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/pkg/requestid"
	"github.com/xiaowuzai/payroll/internal/pkg/response"
	"github.com/xiaowuzai/payroll/internal/service"
	"net/http"
)

type UserHandler struct {
	user   *service.UserService
	logger *logger.Logger
}

func NewUserHandler(userService *service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		user:   userService,
		logger: logger,
	}
}

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username" binding:"required"`
	AccountName string `json:"accountName" binding:"required"`
	Email       string `json:"email"`
	RoleId      string `json:"roleId" binding:"required"`
	RoleName    string `json:"roleName"`
	Password    string `json:"password" binding:"required"`
	Status      int32  `json:"status"`
	Created     int64  `json:"created"`
}

func (u *User) toService() *service.User {
	return &service.User{
		Id:          u.Id,
		Username:    u.Username,
		AccountName: u.AccountName,
		Email:       u.Email,
		RoleId:      u.RoleId,
		Password:    u.Password,
		Status:      u.Status,
	}
}

func (u *User) fromService(su *service.User) {
	u.Id = su.Id
	u.Username = su.Username
	u.AccountName = su.AccountName
	u.Email = su.Email
	u.RoleId = su.RoleId
	u.RoleName = su.RoleName
	u.Status = su.Status
	u.Created = su.Created.Unix()
}

func (uh *UserHandler) AddUser(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("AddUser function called")

	user := &User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		log.Error("AddUser ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	err = uh.user.AddUser(ctx, &service.User{
		Username:    user.Username,
		AccountName: user.AccountName,
		Email:       user.Email,
		RoleId:      user.RoleId,
		Password:    user.Password,
		Status:      user.Status,
	})
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.Success(c,nil)
}

func (uh *UserHandler) AddAdmin(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("AddAdmin function called")

	err := uh.user.AddUser(ctx, &service.User{
		Username:    "管理员",
		AccountName: "admin",
		Email:       "",
		RoleId:      "4fba1999-a7f9-4b34-b82a-b3f34cfe4d81",
		Password:    "Admin@123",
	})
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("AddAdmin function success")
	c.JSON(http.StatusOK, nil)
}

func (uh *UserHandler) ListUser(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("AddAdmin function called")

	res, err := uh.user.ListUser(ctx)
	if err != nil {
		response.WithError(c, err)
		return
	}

	data := make([]*User, 0, len(res))
	for _, v := range res {
		sUser := v
		user := &User{}
		user.fromService(sUser)

		data = append(data, user)
	}

	response.Success(c, data)
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("GetUser function called")

	id, has := c.Params.Get("id")
	if id == "" || !has {
		response.ParamsError(c, "id不存在")
		return
	}

	su, err := uh.user.GetUser(ctx, id)
	if err != nil {
		response.WithError(c, err)
		return
	}

	user := &User{}
	user.fromService(su)

	log.Info("GetUser function success")
	response.Success(c, user)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("UpdateUser function called")

	user := &User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		log.Error("UpdateUser ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}
	if user.Id == "" {
		log.Error("UpdateUser id not exist")
		response.ParamsError(c, "id不存在")
		return
	}

	su := user.toService()

	err = uh.user.UpdateUser(ctx, su)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("UpdateUser function success")
	response.Success(c,nil)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := uh.logger.WithRequestId(ctx)
	log.Info("DeleteUser function called")

	req := &RequestId{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Error("DeleteUser ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}
	if req.Id == "" {
		log.Error("DeleteUser id not exist")
		response.ParamsError(c, "id不存在")
		return
	}

	err = uh.user.DeleteUser(ctx, req.Id)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("DeleteUser function success")
	response.Success(c, nil)
}
