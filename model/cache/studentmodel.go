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
	cacheStudentIdPrefix        = "cache:student:id:"
	cacheStudentClassNamePrefix = "cache:student:class:name:"
)

type (
	StudentModel interface {
		Insert(data *Student) error
		FindOne(id int64) (*Student, error)
		FindOneByClassName(class string, name string) (*Student, error)
		Update(data *Student) error
		Delete(id int64) error
	}

	defaultStudentModel struct {
		dbConn db.DBConn
		table  string
	}

	Student struct {
		Id         int64           `gorm:"column:id"`
		Class      string          `gorm:"column:class"`
		Name       string          `gorm:"column:name"`
		Age        sql.NullInt64   `gorm:"column:age"`
		Score      sql.NullFloat64 `gorm:"column:score"`
		CreateTime time.Time       `gorm:"column:create_time"`
		UpdateTime sql.NullTime    `gorm:"column:update_time"`
	}
)

func NewStudentModel(gdb *gorm.DB, c cache.CacheConf, opts ...cache.Option) StudentModel {
	return &defaultStudentModel{
		dbConn: db.NewDBConn(gdb, c, opts...),
		table:  "`student`",
	}
}

func (m *defaultStudentModel) Insert(data *Student) error {
	studentIdKey := fmt.Sprintf("%s%v", cacheStudentIdPrefix, data.Id)
	studentClassNameKey := fmt.Sprintf("%s%v:%v", cacheStudentClassNamePrefix, data.Class, data.Name)
	err := m.dbConn.InsertIndex(func(conn *gorm.DB) error {
		return conn.Create(data).Error
	}, studentIdKey, studentClassNameKey)
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

func (m *defaultStudentModel) FindOne(id int64) (*Student, error) {
	var resp Student

	studentIdKey := fmt.Sprintf("%s%v", cacheStudentIdPrefix, id)
	err := m.dbConn.FindIndex(&resp, studentIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("id = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &resp, nil
}

func (m *defaultStudentModel) FindOneByClassName(class string, name string) (*Student, error) {
	var resp Student

	studentClassNameKey := fmt.Sprintf("%s%v:%v", cacheStudentClassNamePrefix, class, name)
	err := m.dbConn.FindIndex(&resp, studentClassNameKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("class,name = ?", class, name).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return &resp, nil
}

func (m *defaultStudentModel) Update(data *Student) error {
	studentIdKey := fmt.Sprintf("%s%v", cacheStudentIdPrefix, data.Id)
	studentClassNameKey := fmt.Sprintf("%s%v:%v", cacheStudentClassNamePrefix, data.Class, data.Name)
	err := m.dbConn.UpdateIndex(func(conn *gorm.DB) error {
		return conn.Updates(data).Error
	}, studentIdKey, studentClassNameKey)
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

func (m *defaultStudentModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	studentClassNameKey := fmt.Sprintf("%s%v:%v", cacheStudentClassNamePrefix, data.Class, data.Name)
	studentIdKey := fmt.Sprintf("%s%v", cacheStudentIdPrefix, id)
	err = m.dbConn.DelIndex(func(conn *gorm.DB) error {
		return conn.Delete(&Student{}, id).Error
	}, studentIdKey, studentClassNameKey)
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

func (m *defaultStudentModel) FormatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheStudentIdPrefix, primary)
}
