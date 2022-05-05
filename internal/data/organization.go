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
	Name         string    `xorm:"name varchar(255)  unique"`
	Path         string    `xorm:"path varchar(255) notnull"` // .ParentId.Id
	//AliasName    string    `xorm:"varchar(36)"`               //   工资类型：手动输入
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
	FeeType      int32     `xorm:"fee_type"` // 1:工资 2:福利 3: 退休    费用类型
	Type         int32     `xorm:"notnull"`  // 1 单位、 2 工资表
	EmployeeType int32     `xorm:"notnull"`  // 员工类型： 1: 公务员  2:事业 3: 企业
}

func (org *Organization) toService() *service.Organization {
	return &service.Organization{
		Id:           org.Id,
		ParentId:     org.ParentId,
		Name:         org.Name,
		Type:         service.OrganizationType(org.Type),
		//AliasName:    org.AliasName,
		FeeType:      service.FeeType(org.FeeType),
		EmployeeType: service.EmployeeType(org.EmployeeType),
	}
}

func (org *Organization) fromService(so *service.Organization) {
	org.Id = so.Id
	org.EmployeeType = int32(so.EmployeeType)
	//org.AliasName = so.AliasName
	org.Type = int32(so.Type)
	org.ParentId = so.ParentId
	org.Name = so.Name
	org.FeeType = int32(so.FeeType)
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
	log := or.logger.WithRequestId(ctx)
	if sorg.ParentId == "" {
		log.Infof("AddOrganization %s : parentId = %s", "父节点错误", sorg.ParentId)
		return errParent
	}

	session := NewSession(ctx, or.data.db)

	// 查看父节点是否存在
	parent := &Organization{
		Id: sorg.ParentId,
	}
	has, err := parent.get(ctx, session, or.logger)
	if err != nil {
		return err
	}
	if !has {
		log.Infof("AddOrganization %s : parentId = %s", "父节点不存在", sorg.ParentId)
		log.Infof("Organization: %+v", parent)
		return errParent
	}
	// 父节点不是组织机构类型，不能添加子节点
	if parent.Type != 1 {
		log.Infof("AddOrganization %s : parentId = %s", "父节点不是组织机构类型", sorg.ParentId)
		return errParent
	}

	org := &Organization{}
	org.fromService(sorg)
	org.Id = uuid.CreateUUID()
	org.Path = parent.Path + "." + org.Id

	session, err = Begin(ctx, session, or.logger)
	if err != nil {
		return err
	}
	err = org.insert(ctx, session, or.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}


	aliasList := org.newAliasList(ctx,  sorg.AliasName)
	alias := PayrollAlias{}
	err = alias.insertList(ctx, session, or.logger, aliasList)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
}

func (or *organizationRepo) UpdateOrganization(ctx context.Context, s *service.Organization) error {
	session, err  := BeginSession(ctx, or.data.db, or.logger)
	if err != nil {
		return err
	}

	org := new(Organization)
	org.fromService(s)

	// 删掉之前的
	alias := &PayrollAlias{OrganizationId: org.Id}
	err = alias.deleteByOrgId(ctx, session, or.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	// 插入新的
	aliasList := org.newAliasList(ctx, s.AliasName)
	err = alias.insertList(ctx,session, or.logger, aliasList)
	if err != nil {
		_ = session.Rollback()
		return err
	}


	err = org.update(ctx, session, or.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_ = session.Commit()
	return nil
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

	alias := &PayrollAlias{OrganizationId: org.Id}
	aliasList, err := alias.listByOrgId(ctx, session, or.logger)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(aliasList))
	for _, v :=  range aliasList  {
		names = append(names, v.AliasName)
	}

	so := org.toService()
	so.AliasName = names

	return so, nil
}

func (or *organizationRepo) DeleteOrganization(ctx context.Context, id string) error {
	session := NewSession(ctx, or.data.db)

	// 查一下
	org := &Organization{Id: id}
	has, err := org.get(ctx, session, or.logger)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	// 如果类型是组织机构
	if service.OrganizationType(org.Type) == service.OType_Organization {
		// 查看是否有子节点
		pOrg := Organization{ParentId: org.Id}
		has, err = pOrg.get(ctx, session, or.logger)
		if err != nil {
			return err
		}
		if has {
			return errors.New("由于存在子节点,该组织机构不能删除")
		}
	}else {  // 如果是工资类型
		// 查看是否有员工属于该工资表
		payInfo := &EmployeeBankInfo{}
		payInfos, err := payInfo.listByOrgId(ctx, session, or.logger, org.Id)
		if err != nil {
			return err
		}
		if len(payInfos) != 0 {
			return errors.New("由于该工资表下有员工，所以不能删除")
		}
	}

	session, err = Begin(ctx, session, or.logger)
	if err != nil {
		return err
	}
	defer session.Close()

	alias := &PayrollAlias{OrganizationId: org.Id}
	err = alias.deleteByOrgId(ctx, session, or.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	deleteOrg := &Organization{Id:id}
	err =deleteOrg.delete(ctx, session, or.logger)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()
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

func (org *Organization) newAliasList(ctx context.Context, names []string) []*PayrollAlias{
	aliasList := make([]*PayrollAlias, 0, len(names))
	for _, name := range names {
		aliasList = append(aliasList, &PayrollAlias{
			Id: uuid.CreateUUID(),
			OrganizationId: org.Id,
			AliasName: 	name,
		})
	}
	return aliasList
}
