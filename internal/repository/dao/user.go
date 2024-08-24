package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrDuplicateUser = errors.New("邮箱已被注册")

type User struct {
	Id       int64          `gorm:"primaryKey"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	// 时间用毫秒数表示
	CreateTime int64
	UpdateTime int64
}

type BasicUserDAO interface {
	Insert(ctx context.Context, user User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) *GORMUserDAO {
	return &GORMUserDAO{db: db}
}

func (u *GORMUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreateTime = now
	user.UpdateTime = now
	err := u.db.WithContext(ctx).Create(&user).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const uniqueConflictsErrNo uint16 = 1062
		if me.Number == uniqueConflictsErrNo {
			return ErrDuplicateUser
		}
	}
	return err
}
