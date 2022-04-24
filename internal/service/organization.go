package service

import (
	"context"
)

type OrganizationRepo interface {
	ListOrganization(context.Context) ([]*Organization, error)
	AddOrganization(context.Context, *Organization) error
	DeleteOrganization(context.Context, string) error
	UpdateOrganization(context.Context, *Organization) error
	GetOrganization(context.Context, string) (*Organization, error)
}

type OrganizationService struct {
	repo OrganizationRepo
}

type Organization struct {
	Id           string          `json:"id"`
	ParentId     string          `json:"parentId"`
	Name         string          `json:"name"`
	SalaryType   string          `json:"salaryType"`   // 工资类型 手动输入。
	FeeType      int32           `json:"feeType"`      // 0:工资 1:福利 2: 退休
	EmployeeType int32           `json:"employeeType"` // 员工类型： 0: 公务员  1:事业 2: 企业
	Type         int32           `json:"type"`         // 0 单位、 1 工资表
	Children     []*Organization `json:"children"`
}

func NewOrganizationService(repo OrganizationRepo) *OrganizationService {
	return &OrganizationService{
		repo: repo,
	}
}

func (os *OrganizationService) ListOrganization(ctx context.Context) (*Organization, error) {
	orgs, err := os.repo.ListOrganization(ctx)
	if err != nil {
		return nil, err
	}

	orgMap := make(map[string][]*Organization, len(orgs))
	var root *Organization
	var rootParent = "root"

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

func buildTree(node *Organization, orgMap map[string][]*Organization) {
	if node == nil || orgMap == nil || len(orgMap) == 0 {
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

func (os *OrganizationService) AddOrganization(ctx context.Context, organization *Organization) error {
	return os.repo.AddOrganization(ctx, organization)
}

func (os *OrganizationService) UpdateOrganization(ctx context.Context, organization *Organization) error {
	return os.repo.UpdateOrganization(ctx, organization)
}

func (os *OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	return os.repo.DeleteOrganization(ctx, id)
}

func (os *OrganizationService) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	return os.repo.GetOrganization(ctx, id)
}
