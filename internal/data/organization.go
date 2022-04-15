package data

import (
	"context"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"log"
	"time"
	"xorm.io/xorm"
)

var _ service.OrganizationRepo = (*organizationRepo)(nil)

type Organization struct {
	Id string `xorm:"id varchar(36) pk"`
	ParentId string `xorm:"parent_id varchar(36) notnull"`
	Name string `xorm:"name varchar(255)  notnull"`
	Path string `xorm:"path varchar(255) notnull"`// .ParentId.Id
	SalaryType string  `xorm:"salary_type varchar(36)"` //   工资类型：手动输入
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	FeeType int32  `xorm:"fee_type int notnull"` // 0:工资 1:福利 2: 退休    费用类型
	Type int32  `xorm:"type int notnull"`  // 0 单位、 1 工资表
	EmployeeType int32 `xorm:"employee_type int notnull"`// 员工类型： 0: 公务员  1:事业 2: 企业
}

func (org *Organization)toService() *service.Organization{
	return &service.Organization{
		Id: org.Id,
		ParentId: org.ParentId,
		Name: org.Name,
		Type: org.Type,
		SalaryType: org.SalaryType,
		FeeType: org.FeeType,
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

type organizationRepo struct {
	data *Data
}

func NewOrganizationRepo(data *Data) service.OrganizationRepo{
	return &organizationRepo{
		data:data,
	}
}

func (or *organizationRepo)ListOrganization(ctx context.Context)([]*service.Organization, error) {
	orgs := make([]*Organization,0)
	err := or.data.db.Find(&orgs)
	if err != nil {
		return nil, err
	}

	res := make([]*service.Organization,0, len(orgs))
	for _, v := range orgs {
		so := v.toService()
		res = append(res, so)
	}

	return res, nil
}

func (or *organizationRepo)AddOrganization(ctx context.Context, sorg *service.Organization) error {
	var errParent = errors.New("父节点错误")
	if sorg.ParentId == "" {
		return errParent
	}

	session, err := BeginSession(ctx, or.data.db)
	if err != nil {
		return err
	}
	defer session.Close()

	// 查看父节点是否存在
	parent := &Organization{
		Id: sorg.ParentId,
	}
	has, err := parent.get(ctx,session)
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
	org.Path = parent.Path+"."+ org.Id


	err = org.insert(ctx, session)
	if err != nil {
		session.Rollback()
		return err
	}

	return nil
}

func (or *organizationRepo)UpdateOrganization(ctx context.Context, s *service.Organization) error {
	session, err := BeginSession(ctx, or.data.db)
	if err != nil {
		return err
	}
	defer session.Close()

	org := new(Organization)
	org.fromService(s)
	err = org.update(ctx, session)
	if err != nil {
		_ = session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

func (or *organizationRepo)GetOrganization(ctx context.Context, id string) (*service.Organization, error) {
	session, err := BeginSession(ctx, or.data.db)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	org := &Organization{Id: id}
	has, err := org.get(ctx, session)
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}
	if !has {
		_ = session.Rollback()
		return nil, errors.New("数据不存在")
	}

	_ = session.Commit()

	return org.toService(), nil
}

func (or *organizationRepo)DeleteOrganization(ctx context.Context, id string) error {

	return nil
}

func (org *Organization) get(ctx context.Context, session *xorm.Session) (bool, error){
	has, err := session.Get(org)
	if err != nil {
		log.Printf("organization get error: %s\n", err.Error())
		return false, errors.New("获取组织信息错误")
	}

	return has, nil
}

func (org *Organization) update(ctx context.Context, session *xorm.Session) error{
	 _, err := session.ID(org.Id).Update(org)
	 if err != nil {
		 log.Printf("organization update error: %s\n", err.Error())
		 return errors.New("更新组织信息错误")
	 }
	 return nil
}

func (org *Organization)insert(ctx context.Context, session *xorm.Session) error {
	_, err := session.Insert(org)
	if err != nil {
		log.Println("organization insert error: ", err.Error())
		return errors.New("数据插入失败")
	}
	return nil
}

func (org *Organization)delete(ctx context.Context, session *xorm.Session) error {
	_, err := session.ID(org.Id).Delete(org)
	if err != nil {
		log.Println("organization insert error: ", err.Error())
		return errors.New("数据删除失败")
	}
	return nil
}
