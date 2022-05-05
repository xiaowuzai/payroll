package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"xorm.io/xorm"
)

type EmployeeBankInfo struct {
	Id             string `xorm:"varchar(36) pk"`
	EmployeeId     string `xorm:"varchar(36) notnull"`
	BankId         string `xorm:"varchar(36) notnull"`
	CardNumber     string `xorm:"varchar(20) notnull unique"`
	OrganizationId string `xorm:"varchar(36) notnull"`
}

func (pi *EmployeeBankInfo) toService() *service.EmployeeBankInfo {
	return &service.EmployeeBankInfo{
		Id:             pi.Id,
		EmployeeId:     pi.EmployeeId,
		BankId:         pi.BankId,
		CardNumber:     pi.CardNumber,
		OrganizationId: pi.OrganizationId,
	}
}

func (pi *EmployeeBankInfo) fromService(spi *service.EmployeeBankInfo) {
	pi.Id = spi.Id
	pi.EmployeeId = spi.EmployeeId
	pi.BankId = spi.BankId
	pi.CardNumber = spi.CardNumber
	pi.OrganizationId = spi.OrganizationId
}

func (pi *EmployeeBankInfo) list(ctx context.Context, session *xorm.Session, logger *logger.Logger, employeeId string) ([]*EmployeeBankInfo, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("EmployeeBankInfo list input employee_id = %s\n", employeeId)

	infos := make([]*EmployeeBankInfo, 0)
	err := session.Where("employee_id = ?", employeeId).Find(&infos)
	if err != nil {
		return nil, errors.ErrDataGet(err.Error())
	}

	return infos, nil
}

func (pi *EmployeeBankInfo) insertList(ctx context.Context, session *xorm.Session, logger *logger.Logger, data []*EmployeeBankInfo) error {
	log := logger.WithRequestId(ctx)
	log.Infof("EmployeeBankInfo insertList input \n")

	_, err := session.Insert(&data)
	if err != nil {
		log.Error("EmployeeBankInfo insertList error: ", err.Error())
		return errors.ErrDataInsert(err.Error())
	}

	return nil
}

func (pi *EmployeeBankInfo) deleteByEmployeeId(ctx context.Context, session *xorm.Session, logger *logger.Logger, empId string) error {
	log := logger.WithRequestId(ctx)
	log.Infof("deleteByEmployeeId input empId = %s\n", empId)

	payInfo := &EmployeeBankInfo{}
	_, err := session.Where("employee_id = ?", empId).Delete(payInfo)
	if err != nil {
		log.Error("deleteByEmployeeId error: ", err.Error())
		return errors.ErrDataDelete(err.Error())
	}

	return nil
}

func (pi *EmployeeBankInfo) listByOrgId(ctx context.Context, session *xorm.Session, logger *logger.Logger, organizationId string) ([]*EmployeeBankInfo, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("EmployeeBankInfo list input employee_id = %s\n", organizationId)

	infos := make([]*EmployeeBankInfo, 0)
	err := session.Where("organization_id = ?", organizationId).Find(&infos)
	if err != nil {
		return nil, errors.ErrDataGet(err.Error())
	}

	return infos, nil
}
