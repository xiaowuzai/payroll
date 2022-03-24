package service

import (
	"context"
	"time"
)

type RoleRepo interface {
	// ctx,userId, *role
	AddRole(context.Context,string, *Role) error
	// ctx,userId
	ListRole(context.Context, string) ([]*Role,error)
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
	Id string
	Description string
	Roles []string
	Created time.Time
}

func NewRoleService(repo RoleRepo) *RoleService{
	return &RoleService{repo: repo}
}

func (r *RoleService)AddRole(ctx context.Context,userId string, role *Role) error {
	return r.repo.AddRole(ctx, userId, role)
}

func  (r *RoleService)ListRole(ctx context.Context, userId string) ([]*Role, error) {
	return  r.repo.ListRole(ctx, userId)
}
