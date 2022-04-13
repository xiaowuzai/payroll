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
	menu *handler.MenuHandler
}

func NewRouter(role *handler.RoleHandler, org *handler.OrganizationHandler,
	user *handler.UserHandler, menu *handler.MenuHandler) *Router {
	return &Router{
		role:role,
		org:org,
		user:user,
		menu:menu,
	}
}

func (r *Router)WithEngine(engine *gin.Engine) {

	engine.Use(middleware.CORSMiddleware(), gin.Recovery(), gin.Logger())
	v1 := engine.Group("/v1")

	v1.POST("/login", r.user.Login)
	v1.POST("/admin-user", r.user.AddAdmin)

	v1auth := v1.Group("/auth")
	//v1auth.Use(middleware.JWTAuthMiddleware())
	v1auth.GET("/whoami", r.user.WhoAmI)

	menu := v1auth.Group("/menu")
	menu.GET("",r.menu.ListMenu)

	role := v1auth.Group("/role")
	role.POST("", r.role.AddRole)
	role.GET("/:id", r.role.GetRole)
	role.GET("", r.role.ListRole)

	organization := v1auth.Group("/organization")
	organization.GET("",r.org.ListOrganization)
	organization.POST("",r.org.AddOrganization)
	//organization.GET("/:id",r.org.ListOrganization)
}
