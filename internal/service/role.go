package service

import (
	"context"
	"time"
)

type RoleRepo interface {
	AddRole(context.Context, *Role) error
}

type RoleService struct {
	repo RoleRepo
}

type Role struct {
	Id string
	Description string
	Roles string
	Created time.Time
}

func NewRoleService(repo RoleRepo) *RoleService{
	return &RoleService{repo: repo}
}

func (r *RoleService)AddRole(ctx context.Context, role *Role) error {
	return r.repo.AddRole(ctx, role)
}

