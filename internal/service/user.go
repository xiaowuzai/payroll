package service

import (
	"context"
	"time"
)

type UserRepo interface {
	AddUser(context.Context, *User)	error
	ChangePasswd(context.Context, string, string, string) error
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
