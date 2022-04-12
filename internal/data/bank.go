package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/service"
	"time"
)

var _ service.BankRepo = (*BankRepo)(nil)

type Bank struct {
	Id string `xorm:"id varchar(36) pk notnull"`
	Name string `xorm:"name varchar(36) unique notnull"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type BankRepo struct {
	data *Data
}

func NewBankRepo( data *Data) *BankRepo{
	return &BankRepo{
		data: data,
	}
}


func (br *BankRepo)AddBank(ctx context.Context, sbank *service.Bank) error {
	bank := &Bank{
		Name: sbank.Name,
	}
	_, err := br.data.db.Insert(bank)
	return err
}

func (br *BankRepo)GetBankList(ctx context.Context) ([]*service.Bank, error) {
	banks := make([]*Bank, 0)
	err := br.data.db.Find(&banks)
	if err != nil {
		return nil, err
	}

	result := make([]*service.Bank, 0, len(banks))
	for _, bank := range banks {
		result = append(result, bank.toService())
	}
	return result, nil
}

func (b *Bank)toService() *service.Bank {
	return &service.Bank{
		Id: b.Id,
		Name: b.Name,
	}
}

func (b *Bank)fromService(sBank *service.Bank ) {
	b.Id= sBank.Id
	b.Name= sBank.Name
}
