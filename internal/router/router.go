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
	org  *handler.OrganizationHandler
	user *handler.UserHandler
	menu *handler.MenuHandler
	employee *handler.EmployeeHandler
}

func NewRouter(role *handler.RoleHandler, org *handler.OrganizationHandler,
	user *handler.UserHandler, menu *handler.MenuHandler) *Router {
	return &Router{
		role: role,
		org:  org,
		user: user,
		menu: menu,
	}
}

func (r *Router) WithEngine(engine *gin.Engine) {

	engine.Use(middleware.CORSMiddleware(), gin.Recovery(), gin.Logger())
	v1 := engine.Group("/v1")

	v1.POST("/login", r.user.Login)
	v1.POST("/admin-user", r.user.AddAdmin)

	v1auth := v1.Group("/auth")
	v1auth.Use(middleware.JWTAuthMiddleware())
	// /v1/auth/whoami
	v1auth.GET("/whoami", r.user.WhoAmI)

	//  /v1/auth/menu
	menu := v1auth.Group("/menu")
	menu.GET("", r.menu.ListMenu)

	//  /v1/auth/role
	role := v1auth.Group("/role")
	role.POST("", r.role.AddRole)
	role.GET("/:id", r.role.GetRole)
	role.GET("", r.role.ListRole)
	role.PUT("", r.role.UpdateRole)
	role.DELETE("", r.role.DeleteRole)

	// /v1/auth/organization
	organization := v1auth.Group("/organization")
	organization.GET("", r.org.ListOrganization)
	organization.POST("", r.org.AddOrganization)
	organization.PUT("", r.org.UpdateOrganization)
	organization.GET("/:id", r.org.GetOrganization)
	organization.DELETE("/:id", r.org.DeleteOrganization)

	// /v1/auth/user
	user := v1auth.Group("/user")
	user.GET("", r.user.ListUser)
	user.GET("/:id", r.user.GetUser)
	user.POST("", r.user.AddUser)
	user.PUT("", r.user.UpdateUser)
	user.DELETE("", r.user.DeleteUser)

	// /v1/auth/employee
	employee := v1auth.Group("/employee")
	employee.GET("", r.employee.ListEmployee)
	employee.GET("/:id", r.employee.GetEmployee)
	employee.POST("", r.employee.AddEmployee)
	employee.PUT("", r.employee.UpdateEmployee)
	employee.DELETE("", r.employee.DeleteEmployee)
}
