package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/service"
	"net/http"
	"time"
)

type RoleHandler struct {
	role service.RoleService
}

func NewRoleHandler(role service.RoleService) *RoleHandler{
	return &RoleHandler{
		role:role,
	}
}

type Role struct {
	Id string `json:"id"`
	Description string `json:"description"`
	Roles string `json:"roles"`  //eg: key1.key2.key3
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

	err = r.role.AddRole(ctx, &service.Role{
		Description: role.Description,
		Roles: role.Roles,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest,err)
		return
	}

	c.JSON(http.StatusOK,nil)
}