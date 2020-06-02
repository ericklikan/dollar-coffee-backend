package test

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/require"
)

func SetupTestDb(t *testing.T) *gorm.DB {
	dbConn, err := gorm.Open("postgres", "postgresql://dev:devpassword@localhost:5432/dollarcoffee?sslmode=disable")
	require.NoError(t, err)
	return dbConn
}
