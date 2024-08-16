package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlErr "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestUserDAO_Insert(t *testing.T) {
	now := time.Now().UnixMilli()
	testCases := []struct {
		name        string
		mock        func(t *testing.T) *sql.DB
		ctx         context.Context
		u           User
		expectedErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				result := sqlmock.NewResult(int64(1), 1)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(result)
				return db
			},
			ctx: context.Background(),
			u: User{
				Email:      sql.NullString{String: "123@qq.com", Valid: true},
				Password:   "123456",
				CreateTime: now,
				UpdateTime: now,
			},
			expectedErr: nil,
		},
		{
			name: "邮箱已存在",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(&mysqlErr.MySQLError{Number: uint16(1062)})
				return db
			},
			ctx: context.Background(),
			u: User{
				Email:      sql.NullString{String: "123@qq.com", Valid: true},
				Password:   "123456",
				CreateTime: now,
				UpdateTime: now,
			},
			expectedErr: ErrDuplicateUser,
		},
		{
			name: "数据库系统出错",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(errors.New("随意一个错误"))
				return db
			},
			ctx: context.Background(),
			u: User{
				Email:      sql.NullString{String: "123@qq.com", Valid: true},
				Password:   "123456",
				CreateTime: now,
				UpdateTime: now,
			},
			expectedErr: errors.New("随意一个错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conn := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				// 连接上替身
				Conn: conn,
				// 跳过show版本（假的，没有版本，得跳过）
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			require.NoError(t, err)
			u := NewGORMUserDAO(db)
			err = u.Insert(tc.ctx, tc.u)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
