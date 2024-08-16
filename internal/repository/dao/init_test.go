package dao

import (
	"fmt"
	"github.com/lyydsheep/awesome-bluebook/config"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestInitTable(t *testing.T) {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	require.NoError(t, err)
	err = InitTable(db)
	require.NoError(t, err)
	fmt.Println()
}
