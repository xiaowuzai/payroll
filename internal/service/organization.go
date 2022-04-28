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

// 工资表类型
type FeeType int32

const (
	FType_PayRoll FeeType = 1 // 工资
	FType_Boon    FeeType = 2 // 福利
	FType_Retire  FeeType = 2 // 福利
)

// 员工类型
type EmployeeType int32

const (
	EType_CivilServant EmployeeType = 1 // 公务员
	EType_Institution  EmployeeType = 2 // 事业单位人员
	EType_Enterprise   EmployeeType = 3 // 企业人员
)

// 组织机构类型
type OrganizationType int32

const (
	OType_Organization OrganizationType = 1
	OType_Payroll      OrganizationType = 2
)

type Organization struct {
	Id           string           `json:"id"`
	ParentId     string           `json:"parentId"`
	Name         string           `json:"name"`
	AliasName    string           `json:"aliasName"`    // 工资类型，手动输入，限制50个汉字。
	FeeType      FeeType          `json:"feeType"`      // 1:工资 2:福利 3: 退休
	EmployeeType EmployeeType     `json:"employeeType"` // 工资性质： 1: 公务员  2:事业 3: 企业
	Type         OrganizationType `json:"type"`         // 1 单位、 2 工资表
	Children     []*Organization  `json:"children"`
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
