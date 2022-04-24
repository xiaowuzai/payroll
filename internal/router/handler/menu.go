package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/pkg/requestid"
	"github.com/xiaowuzai/payroll/internal/pkg/response"
	"github.com/xiaowuzai/payroll/internal/service"
)

type MenuHandler struct {
	menu   *service.MenuService
	logger *logger.Logger
}

func NewMenuHandler(menu *service.MenuService, logger *logger.Logger) *MenuHandler {
	return &MenuHandler{
		menu:   menu,
		logger: logger,
	}
}

type Menu struct {
	MenuKeys map[string]string `json:"menuKeys"`
}

func (mh *MenuHandler) ListMenu(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := mh.logger.WithRequestId(ctx)
	log.Info("ListMenu function called")

	res, err := mh.menu.ListMenu(ctx)
	if err != nil {
		response.WithError(c, err)
		return
	}

	data := Menu{
		MenuKeys: res.MenuKeys,
	}

	log.Info("ListMenu function success")
	response.Success(c, data)
}
