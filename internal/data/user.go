package data

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/xiaowuzai/payroll/internal/pkg/errors"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/service"
	"github.com/xiaowuzai/payroll/pkg/uuid"
	"time"
	"xorm.io/xorm"
)

var _ service.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewUserRepo(data *Data, logger *logger.Logger) service.UserRepo {
	return &userRepo{
		data:   data,
		logger: logger,
	}
}

func (ur *userRepo) AddUser(ctx context.Context, su *service.User) error {
	session := NewSession(ctx, ur.data.db)

	passwd, salt := addSalt(su.Password)

	u := &User{}
	u.fromService(su)
	u.Password = passwd
	u.Salt = salt

	return u.insert(ctx, session, ur.logger)
}

// userId 当前登录的帐号。
// ownerId 被修改者的帐号
func (ur *userRepo) ChangePasswd(ctx context.Context, userId, ownerId, newPasswd string) error {
	session := NewSession(ctx, ur.data.db)
	log := ur.logger.WithRequestId(ctx)

	// 不是管理员，不能修改密码
	user := &User{Id: userId}
	has, err := user.get(ctx, session, ur.logger)
	if err != nil {
		return err
	}
	if !has {
		message := fmt.Sprintf("未找到当前登录的用户, userId = %s", userId)
		log.Error(message)
		return errors.DataNotFound(message)
	}

	// 查看 ownerId 是否存在
	owner := &User{Id: ownerId}
	has, err = owner.get(ctx, session, ur.logger)
	if err != nil {
		return err
	}
	if !has {
		message := fmt.Sprintf("未找到修改密码的用户, userId = %s", ownerId)
		log.Error(message)
		return errors.DataNotFound(message)
	}

	passwd, salt := addSalt(newPasswd)
	owner.Password = passwd
	owner.Salt = salt
	return owner.update(ctx, session, ur.logger)
}

// 认证成功，返回 menu Ids
func (ur *userRepo) Login(ctx context.Context, accountName, passwd string) (*service.User, []string, error) {
	user := &User{AccountName: accountName}
	session := NewSession(ctx, ur.data.db)

	has, err := user.get(ctx, session, ur.logger)
	if err != nil {
		return nil, nil, err
	}
	if !has {
		ur.logger.Error("Login 帐号不存在: accountName = ", accountName)
		return nil, nil, errors.DataNotFound("帐号不存在")
	}

	// 验证密码
	if !validateUserPassword(append([]byte(passwd), user.Salt...), user.Password) {
		return nil, nil, errors.New("密码错误")
	}

	// 获取用户的 menuId列表
	roleMenu := &RoleMenu{}
	roleMenus, err := roleMenu.list(ctx, session, ur.logger, user.RoleId)
	if err != nil {
		return nil, nil, err
	}
	menuIds := make([]string, 0, len(roleMenus))
	for _, v := range roleMenus {
		menuIds = append(menuIds, v.MenuId)
	}

	sUser := &service.User{
		Id:          user.Id,
		Username:    user.Username,
		AccountName: user.AccountName,
		Created:     user.Created,
	}

	return sUser, menuIds, nil
}

func (ur *userRepo) ListUser(ctx context.Context) ([]*service.User, error) {
	session, err := BeginSession(ctx, ur.data.db, ur.logger)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// 获取角色信息
	role := &Role{}
	roles, err := role.list(ctx, session, ur.logger)
	if err != nil {
		return nil, err
	}
	roleMap := make(map[string]string, len(roles))
	for _, r := range roles {
		roleMap[r.Id] = r.Name
	}

	// 获取用户信息
	user := &User{}
	users, err := user.list(ctx, session, ur.logger)
	if err != nil {
		return nil, err
	}

	susers := make([]*service.User, 0, len(users))
	for _, v := range users {
		susers = append(susers, v.toService(roleMap))
	}

	_ = session.Commit()
	return susers, nil
}

func (ur *userRepo) UpdateUser(ctx context.Context, su *service.User) error {
	session := NewSession(ctx, ur.data.db)

	var (
		passwd []byte
		salt   []byte
	)
	if su.Password != "" {
		passwd, salt = addSalt(su.Password)
	}

	u := &User{}
	u.fromService(su)
	u.Password = passwd
	u.Salt = salt

	return u.update(ctx, session, ur.logger)
}

func (ur *userRepo) DeleteUser(ctx context.Context, id string) error {
	session := NewSession(ctx, ur.data.db)

	user := &User{Id: id}
	return user.delete(ctx, session, ur.logger)
}

func (ur *userRepo) GetUser(ctx context.Context, id string) (*service.User, error) {
	session := NewSession(ctx, ur.data.db)
	log := ur.logger.WithRequestId(ctx)

	user := &User{Id: id}
	has, err := user.get(ctx, session, ur.logger)
	if err != nil {
		return nil, err
	}
	if !has {
		message := fmt.Sprintf("User id = %s not exist", id)
		log.Error(message)
		return nil, errors.DataNotFound(message)
	}

	return user.toService(nil), nil
}

type User struct {
	Id          string    `xorm:"id pk varchar(36) notnull"`
	Username    string    `xorm:"username varchar(45) unique"`     // 用户名
	AccountName string    `xorm:"account_name varchar(20) unique"` // 用户帐号，用于登录
	Email       string    `xorm:"email varchar(30)"`
	RoleId      string    `xorm:"role_id varchar(36) notnull"` // Role.Id
	Password    []byte    `xorm:"password varchar(20)"`        // default: 123456
	Salt        []byte    `xorm:"salt varchar(255)"`
	Status      int32     `xorm:"status"` // 0 正常、1 禁用
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}

// roleMap := map[roleId]roleName
func (u *User) toService(roleMap map[string]string) *service.User {
	var (
		roleName string
		ok       bool
	)

	if roleMap != nil {
		roleName, ok = roleMap[u.RoleId]
	}
	_ = ok

	return &service.User{
		Id:          u.Id,
		Username:    u.Username,
		AccountName: u.AccountName,
		Email:       u.Email,
		RoleId:      u.RoleId,
		RoleName:    roleName,
		Status:      u.Status,
		Created:     u.Created,
	}
}

func (u *User) fromService(su *service.User) {
	u.Id = su.Id
	u.Username = su.Username
	u.AccountName = su.AccountName
	u.Email = su.Email
	u.RoleId = su.RoleId
	u.Status = su.Status
}

func (u *User) update(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Info("User update input %+v\n", *u)

	_, err := session.ID(u.Id).Update(u)
	if err != nil {
		log.Error("User update error: ", err.Error())
		return errors.ErrDataUpdate(err.Error())
	}
	return nil
}

func (u *User) delete(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Info("User delete input %+v\n", *u)

	_, err := session.ID(u.Id).Delete(u)
	if err != nil {
		log.Error("User delete error:  ", err.Error())
		return errors.ErrDataDelete(err.Error())
	}
	return nil
}

func (u *User) insert(ctx context.Context, session *xorm.Session, logger *logger.Logger) error {
	log := logger.WithRequestId(ctx)
	log.Info("User insert input %+v\n", *u)

	_, err := session.Insert(u)
	if err != nil {
		log.Error("User insert error: ", err.Error())
		return errors.ErrDataInsert(err.Error())
	}
	return err
}

func (u *User) get(ctx context.Context, session *xorm.Session, logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Infof("User get input %+v\n", *u)

	has, err := session.Get(u)
	if err != nil {
		log.Error("User get error: ", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}

	return has, nil
}

func (u *User) list(ctx context.Context, session *xorm.Session, logger *logger.Logger) ([]*User, error) {
	log := logger.WithRequestId(ctx)
	log.Info("User list input")

	users := make([]*User, 0)
	err := session.Find(&users)
	if err != nil {
		log.Error("User list error: ", err.Error())
		return nil, errors.ErrDataGet(err.Error())
	}

	return users, nil
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
func (u *User) hasByRoleId(ctx context.Context, session *xorm.Session, logger *logger.Logger) (bool, error) {
	log := logger.WithRequestId(ctx)
	log.Info("User insert input %+v\n", *u)

	user := &User{}
	has, err := session.ID(u.Id).Exist(user)
	if err != nil {
		log.Error("User getByRoleId error", err.Error())
		return false, errors.ErrDataGet(err.Error())
	}

	return has, err
}
