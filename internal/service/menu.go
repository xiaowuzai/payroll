package service

import "context"

type MenuRepo interface {
	ListMenu(context.Context)(*Menu, error)
}

type MenuService struct {
	repo MenuRepo
}

func NewMenuService(repo MenuRepo) *MenuService{
	return &MenuService{
		repo: repo,
	}
}

type Menu struct {
	MenuKeys map[string]string
}

func (ms *MenuService) ListMenu(ctx context.Context)(*Menu, error) {
	return ms.repo.ListMenu(ctx)
}





