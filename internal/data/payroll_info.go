package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"xorm.io/xorm"
)

type PayrollInfo struct {
	Id             string `json:"id"`
	EmployeeId     string `json:"employeeId"`
	BankId         string `json:"bankId"`
	CardNumber     string `json:"cardNumber"`
	OrganizationId string `json:"organizationId"`
}

func (pi *PayrollInfo) toService() *service.PayrollInfo {
	return &service.PayrollInfo{
		Id:             pi.Id,
		EmployeeId:     pi.EmployeeId,
		BankId:         pi.BankId,
		CardNumber:     pi.CardNumber,
		OrganizationId: pi.OrganizationId,
	}
}

func (pi *PayrollInfo) fromService(spi *service.PayrollInfo) {
	pi.Id = spi.Id
	pi.EmployeeId = spi.EmployeeId
	pi.BankId = spi.BankId
	pi.CardNumber = spi.CardNumber
	pi.OrganizationId = spi.OrganizationId
}

func (pi *PayrollInfo) list(ctx context.Context, session *xorm.Session, logger *logger.Logger, employeeId string) ([]*PayrollInfo, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("PayrollInfo list input employee_id = %s\n", employeeId)

	infos := make([]*PayrollInfo, 0)
	err := session.Where("employee_id = ?", employeeId).Find(&infos)
	if err != nil {
		return nil, errors.ErrDataGet(err.Error())
	}

	return infos, nil
}

func (pi *PayrollInfo) insertList(ctx context.Context, session *xorm.Session, logger *logger.Logger, data []*PayrollInfo) error {
	log := logger.WithRequestId(ctx)
	log.Infof("PayrollInfo insertList input \n")

	_, err := session.Insert(&data)
	if err != nil {
		log.Error("PayrollInfo insertList error: ", err.Error())
		return errors.ErrDataInsert(err.Error())
	}

	return nil
}

func (pi *PayrollInfo) deleteByEmployeeId(ctx context.Context, session *xorm.Session, logger *logger.Logger, empId string) error {
	log := logger.WithRequestId(ctx)
	log.Infof("deleteByEmployeeId input empId = %s\n", empId)

	payInfo := &PayrollInfo{}
	_, err := session.Where("employee_id = ?", empId).Delete(payInfo)
	if err != nil {
		log.Error("deleteByEmployeeId error: ", err.Error())
		return errors.ErrDataDelete(err.Error())
	}

	return nil
}
