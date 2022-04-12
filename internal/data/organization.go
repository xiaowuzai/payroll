package data

import (
	"context"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"time"
)

var _ service.OrganizationRepo = (*organizationRepo)(nil)

// 组织结构中的组织类型
type OrganizationType int32
const (
	OrgTypeEnterprise OrganizationType = iota  // 单位
	OrgTypeSalaryTable  // 工资表
)

// 工资表中的工资类型
type SalaryType int32
const (
	SalTypeSalary SalaryType = iota
	SalTypeWelfare
	SalTypeRetire
)

type EmployeeType int32
const (
	EepTypeCivilServant EmployeeType = iota   // 公务员
	EepTypePublic // 事业单位
	EepTypeEnterprise // 企业单位
)

type Organization struct {
	Id string `xorm:"id varchar(36) pk"`
	ParentId string `xorm:"parent_id varchar(36) notnull"`
	Name string `xorm:"name varchar(255)  notnull"`
	Type OrganizationType  `xorm:"type int notnull"`  // 0 单位、 1 工资表
	Path string `xorm:"path varchar(255) notnull"`// .ParentId.Id
	SalaryType SalaryType  `xorm:"salary_type int notnull"` // 0:工资 1:福利 2: 退休    工资类型
	EmployeeType EmployeeType `xorm:"employee_type int notnull"`// 员工类型： 0: 公务员  1:事业 2: 企业
	Created time.Time `xorm:""`
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
		res = append(res, &service.Organization{
			Id: v.Id,
			Name: v.Name,
			Type: int32(v.Type),   // 0 单位、 1 工资表
			ParentId: v.ParentId,
			SalaryType: int32(v.SalaryType),   // 0:工资 1:福利 2: 退休    工资类型
			EmployeeType: int32(v.EmployeeType), // 员工类型： 0: 公务员  1:事业 2: 企业
		})
	}

	return res, nil
}

func (or *organizationRepo)AddOrganization(ctx context.Context, sorg *service.Organization) error {
	var errParent = errors.New("父节点错误")
	if sorg.ParentId == "" {
		return errParent
	}

	// 查看父节点是否存在
	parentOrg, err := or.getOrganization(ctx, sorg.ParentId)
	if err != nil {
		return err
	}
	if parentOrg == nil {
		return errParent
	}

	org := &Organization{
		Id: uuid.CreateUUID(),
		ParentId: sorg.ParentId,
		Name: sorg.Name,
		Type:OrganizationType(sorg.Type),   // 0 单位、 1 工资表
		SalaryType: SalaryType(sorg.SalaryType),
		EmployeeType: EmployeeType(sorg.EmployeeType),
	}
	org.Path = parentOrg.Path+"."+ org.Id

	err = or.insertOrganization(ctx, org)
	if err != nil {
		return err
	}

	return nil
}

func (or *organizationRepo)getOrganization(ctx context.Context, orgId string)(*Organization, error) {
	org := &Organization{}
	has, err := or.data.db.ID(orgId).Get(org)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil,nil
	}
	return org, nil
}

func (or *organizationRepo)insertOrganization(ctx context.Context, organization *Organization) error{
	num, err := or.data.db.Insert(organization)
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("插入组织失败")
	}
	return nil
}

