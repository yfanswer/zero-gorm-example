package main

import (
	perrors "github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func GDB() *gorm.DB {
	return db
}

func NewGDB(dsn string) (err error) {
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "xy_",
			SingularTable: true,
		},
	})
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}
