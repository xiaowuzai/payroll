package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/router/handler/response"
	"github.com/xiaowuzai/payroll/internal/service"
)

type BankHandler struct {
	bank *service.BankService
}

type Bank struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

func NewBankHandler(bank *service.BankService) *BankHandler{
	return &BankHandler{
		bank: bank,
	}
}

func (bh *BankHandler)AddBank(c *gin.Context) {
	req := &Bank{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.BadRequest(c, err.Error())
	}

	ctx := c.Request.Context()
	bh.bank.AddBank(ctx, &service.Bank{
		Name: req.Name,
	})

}
