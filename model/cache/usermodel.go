package cache

import (
	"database/sql"
	"fmt"
	"github.com/yfanswer/zero-gorm/db"
	"time"

	perrors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
)

var (
	cacheUserIdPrefix     = "cache:user:id:"
	cacheUserMobilePrefix = "cache:user:mobile:"
	cacheUserNamePrefix   = "cache:user:name:"
	cacheUserTpPrefix     = "cache:user:tp:"
	cacheUserUserPrefix   = "cache:user:user:"
)

type (
	UserModel interface {
		Insert(data *User) error
		FindOne(id int64) (*User, error)
		FindOneByMobile(mobile string) (*User, error)
		FindOneByName(name sql.NullString) (*User, error)
		FindOneByTp(tp int64) (*User, error)
		FindOneByUser(user string) (*User, error)
		Update(data *User) error
		Delete(id int64) error
	}

	defaultUserModel struct {
		dbConn db.DBConn
		table  string
	}

	User struct {
		Id         int64          `gorm:"column:id"`
		User       string         `gorm:"column:user"`     // 用户
		Name       sql.NullString `gorm:"column:name"`     // 用户\t名称
		Password   string         `gorm:"column:password"` // 用户密码
		Mobile     string         `gorm:"column:mobile"`   // 手机号
		Gender     string         `gorm:"column:gender"`   // 男｜女｜未公开
		Nickname   string         `gorm:"column:nickname"` // 用户昵称
		Tp         int64          `gorm:"column:type"`     // 用户类型
		CreateTime sql.NullTime   `gorm:"column:create_time"`
		UpdateTime time.Time      `gorm:"column:update_time"`
	}
)

func NewUserModel(gdb *gorm.DB, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &defaultUserModel{
		dbConn: db.NewDBConn(gdb, c, opts...),
		table:  "`user`",
	}
}

func (m *defaultUserModel) Insert(data *User) error {
	userTpKey := fmt.Sprintf("%s%v", cacheUserTpPrefix, data.Tp)
	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, data.User)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	err := m.dbConn.InsertIndex(func(conn *gorm.DB) error {
		return conn.Create(data).Error
	}, userNameKey, userTpKey, userUserKey, userIdKey, userMobileKey)
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

func (m *defaultUserModel) FindOne(id int64) (*User, error) {
	var resp User

	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	err := m.dbConn.FindIndex(&resp, userIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("id = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &resp, nil
}

func (m *defaultUserModel) FindOneByMobile(mobile string) (*User, error) {
	var resp User

	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
	err := m.dbConn.FindIndex(&resp, userMobileKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("mobile = ?", mobile).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &resp, nil
}

func (m *defaultUserModel) FindOneByName(name sql.NullString) (*User, error) {
	var resp User

	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, name)
	err := m.dbConn.FindIndex(&resp, userNameKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("name = ?", name).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &resp, nil
}

func (m *defaultUserModel) FindOneByTp(tp int64) (*User, error) {
	var resp User

	userTpKey := fmt.Sprintf("%s%v", cacheUserTpPrefix, tp)
	err := m.dbConn.FindIndex(&resp, userTpKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("tp = ?", tp).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &resp, nil
}

func (m *defaultUserModel) FindOneByUser(user string) (*User, error) {
	var resp User

	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, user)
	err := m.dbConn.FindIndex(&resp, userUserKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("user = ?", user).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &resp, nil
}

func (m *defaultUserModel) Update(data *User) error {
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userTpKey := fmt.Sprintf("%s%v", cacheUserTpPrefix, data.Tp)
	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, data.User)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	err := m.dbConn.UpdateIndex(func(conn *gorm.DB) error {
		return conn.Updates(data).Error
	}, userTpKey, userUserKey, userIdKey, userMobileKey, userNameKey)
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

func (m *defaultUserModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userTpKey := fmt.Sprintf("%s%v", cacheUserTpPrefix, data.Tp)
	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, data.User)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	err = m.dbConn.DelIndex(func(conn *gorm.DB) error {
		return conn.Delete(&User{}, id).Error
	}, userIdKey, userMobileKey, userNameKey, userTpKey, userUserKey)
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

func (m *defaultUserModel) FormatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}
