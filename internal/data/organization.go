package data

import (
	"context"
	"fmt"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"time"
	"xorm.io/xorm"
)

var _ service.OrganizationRepo = (*organizationRepo)(nil)

type organizationRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewOrganizationRepo(data *Data, logger *logger.Logger) service.OrganizationRepo {
	return &organizationRepo{
		data:   data,
		logger: logger,
	}
}

type Organization struct {
	Id           string    `xorm:"id varchar(36) pk"`
	ParentId     string    `xorm:"parent_id varchar(36) notnull"`
	Name         string    `xorm:"name varchar(255)  notnull"`
	Path         string    `xorm:"path varchar(255) notnull"` // .ParentId.Id
	SalaryType   string    `xorm:"salary_type varchar(36)"`   //   工资类型：手动输入
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
	FeeType      int32     `xorm:"fee_type int notnull"`      // 0:工资 1:福利 2: 退休    费用类型
	Type         int32     `xorm:"type int notnull"`          // 0 单位、 1 工资表
	EmployeeType int32     `xorm:"employee_type int notnull"` // 员工类型： 0: 公务员  1:事业 2: 企业
}

func (org *Organization) toService() *service.Organization {
	return &service.Organization{
		Id:           org.Id,
		ParentId:     org.ParentId,
		Name:         org.Name,
		Type:         org.Type,
		SalaryType:   org.SalaryType,
		FeeType:      org.FeeType,
		EmployeeType: org.EmployeeType,
	}
}

func (org *Organization) fromService(so *service.Organization) {
	org.Id = so.Id
	org.EmployeeType = so.EmployeeType
	org.SalaryType = so.SalaryType
	org.Type = so.Type
	org.ParentId = so.ParentId
	org.Name = so.Name
}

func (or *organizationRepo) ListOrganization(ctx context.Context) ([]*service.Organization, error) {
	orgs := make([]*Organization, 0)
	err := or.data.db.Find(&orgs)
	if err != nil {
		return nil, err
	}

	res := make([]*service.Organization, 0, len(orgs))
	for _, v := range orgs {
		so := v.toService()
		res = append(res, so)
	}

	return res, nil
}

func (or *organizationRepo) AddOrganization(ctx context.Context, sorg *service.Organization) error {
	var errParent = errors.New("父节点错误")
	if sorg.ParentId == "" {
		return errParent
	}

	session, err := BeginSession(ctx, or.data.db, or.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	// 查看父节点是否存在
	parent := &Organization{
		Id: sorg.ParentId,
	}
	has, err := parent.get(ctx, session, or.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	if !has {
		_ = session.Rollback()
		return errParent
	}

	org := &Organization{}
	org.fromService(sorg)
	org.Id = uuid.CreateUUID()
	org.Path = parent.Path + "." + org.Id

	err = org.insert(ctx, session, or.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

func (or *organizationRepo) UpdateOrganization(ctx context.Context, s *service.Organization) error {
	session := NewSession(ctx, or.data.db)

	org := new(Organization)
	org.fromService(s)
	return org.update(ctx, session, or.logger)
}

func (or *organizationRepo) GetOrganization(ctx context.Context, id string) (*service.Organization, error) {
	session := NewSession(ctx, or.data.db)

	org := &Organization{Id: id}
	has, err := org.get(ctx, session, or.logger)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.DataNotFound(fmt.Sprintf("Organization id = %s not found", id))
	}

	return org.toService(), nil
}

func (or *organizationRepo) DeleteOrganization(ctx context.Context, id string) error {

	return nil
}

func (org *Organization) get(ctx context.Context, session *xorm.Session, logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("Organization get input %+v\n", *org)

	has, err := session.Get(org)
	if err != nil {
		log.Errorf("Organization get error: %s\n", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}

	return has, nil
}

func (org *Organization) update(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Organization update input %+v\n", *org)

	_, err := session.ID(org.Id).Update(org)
	if err != nil {
		log.Errorf("Organization update error: %s\n", err.Error())
		return errors.ErrDataUpdate(err.Error())
	}
	return nil
}

func (org *Organization) insert(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Organization insert input %+v\n", *org)

	_, err := session.Insert(org)
	if err != nil {
		log.Errorf("Organization insert error: %s\n", err.Error())
		return errors.ErrDataInsert(err.Error())
	}
	return nil
}

func (org *Organization) delete(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Infof("Organization delete list input %+v\n", *org)

	_, err := session.ID(org.Id).Delete(org)
	if err != nil {
		log.Errorf("Organization delete error: %s\n", err.Error())
		return errors.ErrDataDelete(err.Error())
	}

	return nil
}
