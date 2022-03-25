package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
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
	Name string `json:"name"`
	Description string `json:"description"`
	MenuKeys map[string]string `json:"menu_key"`
	Menus []string `json:"menus"`
	Created time.Time `json:"created"`
}

func (r *RoleHandler) AddRole(c *gin.Context) {
	role := &Role{}
	err := c.ShouldBindJSON(role)
	if err != nil {
		c.JSON(http.StatusBadRequest,err)
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
		c.JSON(http.StatusBadRequest,err)
		return
	}

	c.JSON(http.StatusOK,nil)
}

func (r *RoleHandler) ListRole(c *gin.Context) {
	ctx := c.Request.Context()
	userId := ""
	sRoles, err := r.role.ListRole(ctx,userId)
	if err != nil {
		log.Printf("List Role error: %s\n", err)
		c.JSON(http.StatusBadRequest, err)
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

	c.JSON(http.StatusOK, roles)
}

func (r *RoleHandler) GetRole(c *gin.Context) {
	roleId := c.Query("id")
	if roleId == ""	{
		c.JSON(http.StatusBadRequest, errors.New("role id is empty"))
		return
	}

	ctx := c.Request.Context()
	userId :=  ""
	sRole, err := r.role.GetRole(ctx,userId, roleId)
	if err != nil {
		c.JSON(http.StatusBadGateway, err)
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
