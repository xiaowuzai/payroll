package data

import (
	"context"
	"crypto/sha1"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"log"
	"time"
)

var _ service.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
}

const (
	admin = "admin"
)

func NewUserRepo(data *Data) service.UserRepo {
	return &userRepo{
		data:data,
	}
}


type User struct {
	Id string `xorm:"id pk varchar(36) notnull "`
	Username string `xorm:"username varchar(45) unique 'username_unique'"` // 用户名
	AccountName string `xorm:"account_name varchar(20) unique 'account_name_unique'"`
	Email string `xorm:"email"`
	RoleId string `xorm:"role"` // Role.Id
	Password []byte `xorm:"password"`  // default: 123456
	Salt []byte `xorm:"salt"`
	Status int32 `xorm:"status"`  // 0 正常、1 禁用
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

const defaultPasswd = "123456"

func addSalt(password string) ([]byte, []byte) {
	salt := uuid.CreatUUIDBinary()
	sm := append([]byte(password), salt...)
	h := sha1.New()
	h.Write([]byte(sm))
	return h.Sum(nil), salt
}

func (ur *userRepo)AddUser(ctx context.Context, user *service.User) error {

	passwd, salt := addSalt(user.Password)

	u := &User{
		Id : uuid.CreateUUID(),
		Username: user.Username,
		AccountName: user.AccountName,
		Password: passwd,
		Salt: salt,
		Email: user.Email,
		RoleId: user.RoleId,
		Status: user.Status,
	}

	err := ur.insert(ctx,u)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepo)insert(ctx context.Context, user *User) error {
	_, err :=  ur.data.db.Insert(user)
	return err
}


func (ur *userRepo)ChangePasswd(ctx context.Context, userId, ownerId, newPasswd string) error {
	// 检查是不是管理员
	user, err := ur.getUser(ctx, userId)
	if err != nil {
		return err
	}
	if user.AccountName != admin {
		return errors.New("您不是系统管理员，不能修改密码")
	}

	return ur.changePasswd(ctx, ownerId, newPasswd)
}

func (ur *userRepo)getUser(ctx context.Context, userId string) (*User, error) {
	u := &User{}
	has, err := ur.data.db.ID(userId).Get(u)
	if err != nil {
		return nil, err
	}
	if !has {
		log.Printf("userId: %s not exist\n", userId)
	}
	return u, nil
}

// passwd 为空，则重置为默认密码
func (ur *userRepo)changePasswd(ctx context.Context, userId , passwd string) error {
	if passwd == "" {
		passwd = defaultPasswd
	}
	pwd, salt := addSalt(passwd)
	user := &User{
		Id: userId,
		Password: pwd,
		Salt: salt,
	}
	return ur.update(ctx, user)
}

func (ur *userRepo)update(ctx context.Context, user *User) error{
	_, err := ur.data.db.ID(user.Id).Update(user)
	if err != nil {
		return err
	}
	return nil
}

// 判断是否为管理员
func isAdmin(ctx context.Context) bool {
	return false
}

