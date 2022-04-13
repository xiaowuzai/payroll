package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/router/handler/response"
	"github.com/xiaowuzai/payroll/internal/service"
)

type MenuHandler struct {
	menu *service.MenuService
}

func NewMenuHandler(menu *service.MenuService) *MenuHandler {
	return &MenuHandler{
		menu: menu,
	}
}

type Menu struct {
	MenuKeys map[string]string `json:"menuKeys"`
}

func (mh *MenuHandler) ListMenu(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := mh.menu.ListMenu(ctx)
	if err != nil {
		response.InternalErr(c,err.Error())
		return
	}

	data := Menu{
		MenuKeys: res.MenuKeys,
	}
	response.Success(c, data)
}