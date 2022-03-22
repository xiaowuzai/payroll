package data

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/service"
)

var _ service.OrganizationRepo = (*organizationRepo)(nil)

type Organization struct {
	Id int
	ParentId int
	Name string
	Type int32    // 0 单位、 1 工资表
	Path string // .ParentId.Id
}

type OrganizationSalary  struct {
	Id string   // OrganizationId
	SalaryType int32   // 0:工资 1:福利 2: 退休    工资类型
	EmployeeType int32 // 员工类型： 0: 公务员  1:事业 2: 企业
}

type organizationRepo struct {
	data *Data
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
			Type: v.Type,   // 0 单位、 1 工资表
			ParentId: v.ParentId,
		})
	}

	return res, nil
}

