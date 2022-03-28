package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/router/handler/response"
	"github.com/xiaowuzai/payroll/internal/service"
	"log"
	"net/http"
	"time"
)

type RoleHandler struct {
	role *service.RoleService
}

func NewRoleHandler(role *service.RoleService) *RoleHandler{
	return &RoleHandler{
		role:role,
	}
}

type Role struct {
	Id string `json:"id"`
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
	MenuKeys map[string]string `json:"menu_key"`
	Menus []string `json:"menus" binding:"required"`
	Created time.Time `json:"created"`
}

// @Summary 添加角色
// @Description 指定角色的菜单权限
// @Tags 角色管理
// @Accept application/json
// @Param Role body Role true ""
// @Success 200 {object} Role
// @Router /v1/auth/role [post]
func (r *RoleHandler) AddRole(c *gin.Context) {
	role := &Role{}
	err := c.ShouldBindJSON(role)
	if err != nil {
		response.Error(c,http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()

	userId := ""
	err = r.role.AddRole(ctx,userId, &service.Role{
		Description: role.Description,
		Name: role.Name,
		Menus: role.Menus,
	})
	if err != nil {
		response.Error(c,http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c,nil)
}

// @Summary 获取角色列表
// @Description 根据用户权限获取角色列表
// @Tags 角色管理
// @Accept application/json
// @Success 200 {object} []*Role
// @Router /v1/auth/role [get]
func (r *RoleHandler) ListRole(c *gin.Context) {
	ctx := c.Request.Context()
	userId := ""
	sRoles, err := r.role.ListRole(ctx,userId)
	if err != nil {
		log.Printf("List Role error: %s\n", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	roles := make([]*Role, 0, len(sRoles))
	for _, v := range sRoles {
		roles = append(roles, &Role{
			Id: v.Id,
			Name: v.Name,
			Description: v.Description,
			Created: v.Created,
		})
	}

	response.Success(c, roles)
}

// @Summary 角色管理
// @Description 获取某个角色信息
// @Tags 角色管理
// @Accept application/json
// @Param id query string true "id"
// @Success 200 {object} Role
// @Router /v1/auth/role/{id} [get]
func (r *RoleHandler) GetRole(c *gin.Context) {
	roleId,has := c.Params.Get("id")
	if roleId == "" || !has	{
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	ctx := c.Request.Context()
	userId :=  ""
	log.Println("roleId: ", roleId)
	sRole, err := r.role.GetRole(ctx, userId, roleId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	role := &Role{
		Id: sRole.Id,
		Name: sRole.Name,
		Description: sRole.Description,
		Menus: sRole.Menus,
		MenuKeys: sRole.MenuKey,
	}

	c.JSON(http.StatusOK, role)
}
