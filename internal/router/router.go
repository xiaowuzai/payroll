package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/xiaowuzai/payroll/internal/pkg/middleware"
	"github.com/xiaowuzai/payroll/internal/router/handler"
)

// ProviderSet is router providers.
var ProviderSet = wire.NewSet(NewRouter)

type Router struct {
	role *handler.RoleHandler
	org *handler.OrganizationHandler
	user *handler.UserHandler
}

func (r *Router)WithEngine(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	v1.Use(gin.Recovery(), gin.Logger(),middleware.CORSMiddleware())
	v1.POST("/login", r.user.Login)
	v1.POST("/admin-user", r.user.AddAdmin)

	v1auth := v1.Group("/auth")
	v1auth.Use(middleware.JWTAuthMiddleware())
	v1auth.POST("/role", r.role.AddRole)
	v1auth.GET("/role/:id", r.role.GetRole)
	v1auth.GET("/role", r.role.ListRole)
	v1auth.GET("/organization",r.org.ListOrganization)
}

func NewRouter(role *handler.RoleHandler, org *handler.OrganizationHandler, user *handler.UserHandler) *Router {
	return &Router{
		role:role,
		org:org,
		user:user,
	}
}