package service

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/middleware"
	"time"
)

type UserRepo interface {
	ListUser(context.Context) ([]*User, error)
	AddUser(context.Context, *User) error
	ChangePasswd(context.Context, string, string, string) error
	Login(context.Context, string, string) (*User, []string, error)
	UpdateUser(context.Context, *User) error
	DeleteUser(context.Context, string) error
	GetUser(context.Context, string) (*User, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

type User struct {
	Id          string
	Username    string
	AccountName string
	Email       string
	RoleId      string
	RoleName    string
	Password    string
	Status      int32
	Created     time.Time
}

func (us *UserService) AddUser(ctx context.Context, user *User) error {
	return us.repo.AddUser(ctx, user)
}

func (us *UserService) Login(ctx context.Context, accountName, passwd string) (string, string, error) {
	user, menuIds, err := us.repo.Login(ctx, accountName, passwd)
	if err != nil {
		return "", "", err
	}

	authInfo := &middleware.AuthInfo{
		UId:   user.Id,
		Name:  user.Username,
		Menus: menuIds,
	}

	token, err := middleware.GenerateToken(authInfo)
	if err != nil {
		return "", "", err
	}
	refresh, err := middleware.GenerateRefreshToken(authInfo)
	if err != nil {
		return "", "", err
	}

	return token, refresh, nil
}

func (us *UserService) ListUser(ctx context.Context) ([]*User, error) {
	return us.repo.ListUser(ctx)
}

func (us *UserService) UpdateUser(ctx context.Context, user *User) error {
	return us.repo.UpdateUser(ctx, user)
}

func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	return us.repo.DeleteUser(ctx, id)
}

func (us *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	return us.repo.GetUser(ctx, id)
}
