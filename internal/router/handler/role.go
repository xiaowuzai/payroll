package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/pkg/requestid"
	"github.com/xiaowuzai/payroll/internal/pkg/response"
	"github.com/xiaowuzai/payroll/internal/service"
	"net/http"
)

type RoleHandler struct {
	role   *service.RoleService
	logger *logger.Logger
}

func NewRoleHandler(role *service.RoleService, logger *logger.Logger) *RoleHandler {
	return &RoleHandler{
		role:   role,
		logger: logger,
	}
}

type Role struct {
	Id          string   `json:"id"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Menus       []string `json:"menus" binding:"required"`
	Created     int64    `json:"created"`
}

func (r *Role) toService() *service.Role {
	return &service.Role{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Menus:       r.Menus,
	}
}

func (r *Role) fromService(sr *service.Role) {
	r.Id = sr.Id
	r.Name = sr.Name
	r.Menus = sr.Menus
	r.Created = sr.Created.Unix()
}

// @Summary 添加角色
// @Description 指定角色的菜单权限
// @Tags 角色管理
// @Accept application/json
// @Param Role body Role true ""
// @Success 200 {object} Role
// @Router /v1/auth/role [post]
func (rh *RoleHandler) AddRole(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := rh.logger.WithRequestId(ctx)
	log.Info("AddRole function called")

	role := &Role{}
	err := c.ShouldBindJSON(role)
	if err != nil {
		log.Error("AddRole ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	userId := ""
	sr := role.toService()
	err = rh.role.AddRole(ctx, userId, sr)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("AddRole function success")
	response.Success(c, nil)
}

// @Summary 获取角色列表
// @Description 根据用户权限获取角色列表
// @Tags 角色管理
// @Accept application/json
// @Success 200 {object} []*Role
// @Router /v1/auth/role [get]
func (rh *RoleHandler) ListRole(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := rh.logger.WithRequestId(ctx)
	log.Info("ListRole function called")

	userId := ""
	sRoles, err := rh.role.ListRole(ctx, userId)
	if err != nil {
		response.WithError(c, err)
		return
	}

	roles := make([]*Role, 0, len(sRoles))
	for _, v := range sRoles {
		roles = append(roles, &Role{
			Id:          v.Id,
			Name:        v.Name,
			Menus:       v.Menus,
			Description: v.Description,
			Created:     v.Created.Unix(),
		})
	}

	log.Info("ListRole function success")
	response.Success(c, roles)
}

// @Summary 角色管理
// @Description 获取某个角色信息
// @Tags 角色管理
// @Accept application/json
// @Param id query string true "id"
// @Success 200 {object} Role
// @Router /v1/auth/role/{id} [get]
func (rh *RoleHandler) GetRole(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := rh.logger.WithRequestId(ctx)
	log.Info("GetRole function called")

	roleId, has := c.Params.Get("id")
	if roleId == "" || !has {
		log.Error("GetRole 请求中没有id")
		response.ParamsError(c, "id不存在")
		return
	}

	userId := ""
	sr, err := rh.role.GetRole(ctx, userId, roleId)
	if err != nil {
		response.WithError(c, err)
		return
	}

	role := new(Role)
	role.fromService(sr)

	log.Info("GetRole function success")
	c.JSON(http.StatusOK, role)
}

// @Summary 添加角色
// @Description 指定角色的菜单权限
// @Tags 角色管理
// @Accept application/json
// @Param Role body Role true ""
// @Success 200 {object} Role
// @Router /v1/auth/role [put]
func (rh *RoleHandler) UpdateRole(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := rh.logger.WithRequestId(ctx)
	log.Info("UpdateRole function called")

	role := &Role{}
	err := c.ShouldBindJSON(role)
	if err != nil {
		log.Error("UpdateRole ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	sr := role.toService()
	err = rh.role.UpdateRole(ctx, sr)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("UpdateRole function success")
	response.Success(c, nil)
}

func (rh *RoleHandler) DeleteRole(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := rh.logger.WithRequestId(ctx)
	log.Info("DeleteRole function called")

	req := &RequestId{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Error("DeleteRole ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	if err := rh.role.DeleteRole(ctx, req.Id); err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("DeleteRole function success")
	response.Success(c, nil)
}
