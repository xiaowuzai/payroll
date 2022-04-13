package data

import (
	"bytes"
	"context"
	"crypto/sha1"
	"errors"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"log"
	"time"
	"xorm.io/xorm"
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
	Username string `xorm:"username varchar(45) unique "` // 用户名
	AccountName string `xorm:"account_name varchar(20) unique"` // 用户帐号，用于登录
	Email string `xorm:"email varchar(30)"`
	RoleId string `xorm:"role_id varchar(36) notnull"` // Role.Id
	Password []byte `xorm:"password varchar(20)"`  // default: 123456
	Salt []byte `xorm:"salt varchar(255)"`
	Status int32 `xorm:"status"`  // 0 正常、1 禁用
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

const defaultPasswd = "123456"


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

func (ur *userRepo)ChangePasswd(ctx context.Context, userId, ownerId, newPasswd string) error {
	// 检查是不是管理员
	user, err := ur.getUser(ctx, userId, "")
	if err != nil {
		return err
	}
	if user.AccountName != admin {
		return errors.New("您不是系统管理员，不能修改密码")
	}

	return ur.changePasswd(ctx, ownerId, newPasswd)
}

// 认证成功，返回 menu Ids
func (ur *userRepo) Login(ctx context.Context, accountName, passwd string) (*service.User, []string, error) {
	user, err := ur.getUser(ctx,"", accountName)
	if err != nil {
		return nil,nil, err
	}

	if !validateUserPassword(append([]byte(passwd), user.Salt...), user.Password) {
		return nil, nil, errors.New("密码错误")
	}

	menuIds, err := ur.getMenus(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	sUser := &service.User{
		Id: user.Id,
		Username: user.Username,
		AccountName: user.AccountName,
		Created: user.Created,
	}

	return sUser, menuIds, nil
}

// userId 和 username 取交集。
func (ur *userRepo)getUser(ctx context.Context, userId string, accountName string) (*User, error) {
	u := &User{}

	session := ur.data.db.NewSession()
	defer session.Close()
	if userId != ""  {
		session = session.ID(userId)
	}
	if accountName != "" {
		session = session.Where("account_name = ?", accountName)
	}
	has, err := session.Get(u)
	if err != nil {
		log.Printf("[ getUser ] error: %s", err.Error())
		return nil, errors.New("查询用户信息失败")
	}
	if !has {
		log.Printf("userId: %s not exist\n", userId)
		return nil, errors.New("用户不存在")
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

func (ur *userRepo)insert(ctx context.Context, user *User) error {
	_, err :=  ur.data.db.Insert(user)
	return err
}


func addSalt(password string) ([]byte, []byte) {
	salt := uuid.CreatUUIDBinary()
	sm := append([]byte(password), salt...)
	h := sha1.New()
	h.Write([]byte(sm))
	return h.Sum(nil), salt
}

func validateUserPassword(sumPwd []byte, password []byte) bool {
	h := sha1.New()
	h.Write(sumPwd)
	if bytes.Compare(h.Sum(nil), password) != 0 {
		return false
	}
	return true
}


// 如果不存在返回false
func (u *User) getByRoleId(ctx context.Context, session *xorm.Session) (bool, error) {
	has, err := session.ID(u.Id).Get(u)
	if err != nil {
		log.Println("getByRoleId error", err.Error())
		return false, errors.New("查找用户失败")
	}

	return has, err
}
