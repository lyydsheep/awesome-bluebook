package dao

import (
	"context"
	"gorm.io/gorm"
)

type User struct {
	Id       int64  `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}

var _ UserDAO = (*GORMUserDAO)(nil)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func (dao *GORMUserDAO) Insert(ctx context.Context, user User) error {
	err := dao.db.WithContext(ctx).Create(&user).Error
	return err
}

func NewGORMUserDAO(db *gorm.DB) *GORMUserDAO {
	return &GORMUserDAO{
		db: db,
	}
}
