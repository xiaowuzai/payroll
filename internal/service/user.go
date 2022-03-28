package service

import (
	"context"
	"github.com/xiaowuzai/payroll/internal/pkg/middleware"
	"time"
)

type UserRepo interface {
	AddUser(context.Context, *User)	error
	ChangePasswd(context.Context, string, string, string) error
	Login(context.Context, string, string) (*User, []string, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService{
	return &UserService{
		repo:repo,
	}
}

type User struct{
	Id string
	Username string
	AccountName string
	Email string
	RoleId string
	Password string `xorm:"password"`  // default: 123456
	Status int32 `xorm:"status"`  // 0 正常、1 禁用
	Created time.Time `xorm:"created"`
}

func (us *UserService) AddUser (ctx context.Context, user *User) error{
	return us.repo.AddUser(ctx, user)
}

func (us *UserService) Login(ctx context.Context, accountName, passwd string)(string, string, error)  {
	user, menuIds, err := us.repo.Login(ctx, accountName, passwd)
	if err != nil {
		return "", "", err
	}

	authInfo := &middleware.AuthInfo{
		UId :user.Id,
		Name :user.Username,
		Menus: menuIds,
	}

	token, err := middleware.GenerateToken(authInfo)
	if err != nil {
		return "" , "", err
	}
	refresh, err := middleware.GenerateRefreshToken(authInfo)
	if err != nil {
		return "" , "", err
	}

	return token,refresh, nil
}