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

func (b *Bank)toService() *service.Bank{
	return &service.Bank{
		Id: b.Id,
		Name: b.Name,
	}
}

func (b *Bank) fromService(bs *service.Bank) {
	b.Id = bs.Id
	b.Name = bs.Name
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

	err = bh.bank.AddBank(ctx, req.toService())
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("AddBank function success")
	response.Success(c, nil)
}

func (bh *BankHandler) ListBank(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := bh.logger.WithRequestId(ctx)
	log.Info("ListBank function called")

	sbs, err := bh.bank.GetBankList(ctx)
	if err != nil {
		response.WithError(c, err)
		return
	}

	data := make([]*Bank, 0, len(sbs))
	for _, v := range sbs {
		bank := &Bank{}
		bank.fromService(v)
		data = append(data, bank)
	}

	log.Info("AddBank function success")
	response.Success(c, data)
}
