package service

import (
	"context"
)

type OrganizationRepo interface {
	ListOrganization(context.Context)([]*Organization, error)
}

type OrganizationService struct {
	repo OrganizationRepo
}

type Organization struct {
	Id int
	Name string
	OrganizationSalary OrganizationSalary
	Type int32    // 0 单位、 1 工资表
	ParentId int
	Children []*Organization
}

type OrganizationSalary  struct {
	Id string   // OrganizationId
	SalaryType int32   // 0:工资 1:福利 2: 退休    工资类型
	EmployeeType int32 // 员工类型： 0: 公务员  1:事业 2: 企业
}

func NewOrganizationService(repo OrganizationRepo) *OrganizationService {
	return &OrganizationService{
		repo:repo,
	}
}

func (os *OrganizationService) ListOrganization(ctx context.Context) (*Organization, error) {
	orgs, err  := os.repo.ListOrganization(ctx)
	if err != nil {
		return nil, err
	}

	orgMap := make(map[int][]*Organization, len(orgs))
	var root *Organization
	var rootParent  = 0

	for _, org := range orgs {
		org := org
		if org.ParentId == rootParent {
			root = org
		}

		orgMap[org.ParentId] = append(orgMap[org.ParentId], org)
	}
	
	buildTree(root, orgMap)

	return root, nil
}

func buildTree(node *Organization, orgMap map[int][]*Organization){
	if node == nil || orgMap == nil || len(orgMap) == 0{
		return
	}

	children, ok := orgMap[node.Id]
	if !ok {
		return
	}

	for _, v := range children {
		v := v
		node.Children = append(node.Children, v)
		buildTree(v, orgMap)
	}
}















