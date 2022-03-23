package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/xiaowuzai/payroll/internal/router/handler"
)

// ProviderSet is router providers.
var ProviderSet = wire.NewSet(NewRouter)

type Router struct {
	role *handler.RoleHandler
	org *handler.OrganizationHandler
}

func (r *Router)WithEngine(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	v1.POST("/role", r.role.AddRole)

	v1.GET("/organization",r.org.Organization)
}

func NewRouter(role *handler.RoleHandler, org *handler.OrganizationHandler) *Router {
	return &Router{
		role:role,
		org:org,
	}
}