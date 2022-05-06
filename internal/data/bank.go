package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"time"
	"xorm.io/xorm"
)

var _ service.BankRepo = (*BankRepo)(nil)

type Bank struct {
	Id      string    `xorm:"varchar(36) pk"`
	Name    string    `xorm:"unique"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type BankRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewBankRepo(data *Data, logger *logger.Logger) service.BankRepo {
	return &BankRepo{
		data:   data,
		logger: logger,
	}
}

func (br *BankRepo) AddBank(ctx context.Context, sbank *service.Bank) error {
	session := NewSession(ctx, br.data.db)

	bank := &Bank{
		Name: sbank.Name,
	}
	return bank.insert(ctx, session, br.logger)
}

func (br *BankRepo) GetBankList(ctx context.Context) ([]*service.Bank, error) {
	session := NewSession(ctx, br.data.db)

	bank := &Bank{}
	banks, err := bank.listAll(ctx, session, br.logger)
	if err != nil {
		return nil, err
	}

	result := make([]*service.Bank, 0, len(banks))
	for _, bank := range banks {
		result = append(result, bank.toService())
	}
	return result, nil
}

func (b *Bank) toService() *service.Bank {
	return &service.Bank{
		Id:   b.Id,
		Name: b.Name,
	}
}

func (b *Bank) fromService(sBank *service.Bank) {
	b.Id = sBank.Id
	b.Name = sBank.Name
}

func (b *Bank) insert(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Bank insert %+v\n", *b)

	_, err := session.Insert(b)
	if err != nil {
		log.Errorf("Bank insert error %s\n", err.Error())
		return errors.ErrDataInsert(err.Error())
	}
	return nil
}

func (b *Bank) get(ctx context.Context, session *xorm.Session, logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Bank get %+v\n", *b)

	has, err := session.Get(b)
	if err != nil {
		log.Errorf("Bank get error %s\n", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}
	//if !has {
	//	message := fmt.Sprintf("Bank get not found, bank.id = %s \n", b.Id)
	//	log.Errorf(message)
	//	return false, errors.DataNotFound(message)
	//}

	return has, nil
}

func (b *Bank) listAll(ctx context.Context, session *xorm.Session, logger *logger.Logger) ([]*Bank, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Bank list \n")

	bs := make([]*Bank, 0)
	err := session.Table(&Bank{}).Find(&bs)
	if err != nil {
		log.Errorf("Bank get error %s\n", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}
	return bs, nil
}

func (b *Bank) listByIds(ctx context.Context, session *xorm.Session, logger *logger.Logger, ids []string) ([]*Bank, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Bank list \n")

	bs := make([]*Bank, 0)
	err := session.Table(&Bank{}).In("id", ids).Find(&bs)
	if err != nil {
		log.Errorf("Bank get error %s\n", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}
	return bs, nil
}

func (b *Bank) existByIds(ctx context.Context, session *xorm.Session, logger *logger.Logger, ids []string) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Bank existByIds: %+v", ids)

	banks, err:= b.listByIds(ctx, session, logger , ids)
	if err != nil {
		return false, err
	}

	if len(ids) != len(banks) {
		return false, nil
	}

	return true, nil
}

func (b *Bank) update(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Bank update %+v\n", *b)

	_, err := session.ID(b.Id).Update(b)
	if err != nil {
		log.Errorf("Bank update error %s\n", err.Error())
		return errors.ErrDataUpdate(err.Error())
	}
	return nil
}
