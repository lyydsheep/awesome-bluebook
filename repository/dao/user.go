package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}

var (
	_                UserDAO = (*GORMUserDAO)(nil)
	UserDuplicateErr         = errors.New("邮箱冲突")
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func (dao *GORMUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const uniqueConflictsErrNo uint16 = 1062
		if me.Number == uniqueConflictsErrNo {
			return UserDuplicateErr
		}
	}
	return err
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}
