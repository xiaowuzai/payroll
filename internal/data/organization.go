package data

import (
	"context"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
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

type Organization struct {
	Id int
	ParentId int
	Name string
	Type OrganizationType    // 0 单位、 1 工资表
	Path string // .ParentId.Id
}

type OrganizationSalary  struct {
	Id string   // OrganizationId
	SalaryType SalaryType   // 0:工资 1:福利 2: 退休    工资类型
	EmployeeType int32 // 员工类型： 0: 公务员  1:事业 2: 企业
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
		})
	}

	return res, nil
}

func (or *organizationRepo)AddOrganization(ctx context.Context, sorg service.Organization) error {
	var errParent = errors.New("父节点错误")
	if sorg.ParentId == 0 {
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
		ParentId: sorg.ParentId,
		Name:sorg.Name,
		Type:OrganizationType(sorg.Type),   // 0 单位、 1 工资表

		//TODO: 没想好怎么处理
		Path: parentOrg.Path,
	}

	if org.Type == OrgTypeSalaryTable {
		orgSalary := &OrganizationSalary {
			
		}
	}
}

func (or *organizationRepo)getOrganization(ctx context.Context, orgId int)(*Organization, error) {
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

func (or *organizationRepo)insertOrganizationSalary(ctx context.Context) {

}
