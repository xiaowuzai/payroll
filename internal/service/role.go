package service

import (
	"context"
	"time"
)

type RoleRepo interface {
	// ctx,userId, *role
	AddRole(context.Context, string, *Role) error
	// ctx,userId
	ListRole(context.Context, string) ([]*Role, error)
	// ctx, userId, roleId
	GetRole(context.Context, string, string) (*Role, error)

	UpdateRole(context.Context, *Role) error

	DeleteRole(context.Context, string) error
}

type RoleService struct {
	repo RoleRepo
}

var (
	menuList = []string{
		"salary-management",
		"salary-preview",
		"salary-approve",
		"salary-query",
		"salary-carry-forward",
		"salary-note-history",
		"salary-statistics",
		"cost-statistics",
		"project-cost-statistic",
		"system-config",
		"user-management",
		"role-management",
		"organization",
		"staff-management",
	}
)

type Role struct {
	Id          string
	Name        string
	Description string
	Menus       []string // menu ids
	Created     time.Time
}

func NewRoleService(repo RoleRepo) *RoleService {
	return &RoleService{repo: repo}
}

func (r *RoleService) AddRole(ctx context.Context, userId string, role *Role) error {
	return r.repo.AddRole(ctx, userId, role)
}

func (r *RoleService) ListRole(ctx context.Context, userId string) ([]*Role, error) {
	return r.repo.ListRole(ctx, userId)
}

func (r *RoleService) GetRole(ctx context.Context, userId, roleId string) (*Role, error) {
	return r.repo.GetRole(ctx, userId, roleId)
}

func (r *RoleService) UpdateRole(ctx context.Context, role *Role) error {
	return r.repo.UpdateRole(ctx, role)
}

func (r *RoleService) DeleteRole(ctx context.Context, id string) error {
	return r.repo.DeleteRole(ctx, id)
}
