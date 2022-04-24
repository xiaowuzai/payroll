package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/pkg/requestid"
	"github.com/xiaowuzai/payroll/internal/pkg/response"
	"github.com/xiaowuzai/payroll/internal/service"
)

type BankHandler struct {
	bank   *service.BankService
	logger *logger.Logger
}

type Bank struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewBankHandler(bank *service.BankService, logger *logger.Logger) *BankHandler {
	return &BankHandler{
		bank:   bank,
		logger: logger,
	}
}

func (bh *BankHandler) AddBank(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := bh.logger.WithRequestId(ctx)
	log.Info("AddBank function called")

	req := &Bank{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Error("AddBank ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	err = bh.bank.AddBank(ctx, &service.Bank{
		Name: req.Name,
	})
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("AddBank function success")
	response.Success(c, nil)
}
