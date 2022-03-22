package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/router/handler"
)

type Router struct {
	role *handler.RoleHandler
}

func (r *Router)WithEngine(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	v1.POST("/role", r.role.AddRole)
}