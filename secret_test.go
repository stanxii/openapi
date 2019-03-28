package openapi

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"database/sql"
	"testing"
)

const (
	dbAddr = "testuser:123456@tcp(localhost:3306)/cloudb?charset=utf8&parseTime=True&loc=Local"
)

var (
	sqlKeeper = SqlSecretKeeper{
		getdb(),
		"app",
		"app_key",
		"thekey",
		"app_secret",
	}
)

func TestSqlSecretKeeper_GetSecret(t *testing.T) {
	keeper := SqlSecretKeeper{
		nil,
		"apps",
		"app_key",
		"thekey",
		"app_secret",
	}
	_, err := keeper.GetSecret()
	assert.NotNil(t, err)
	keeper.Db = getdb()
	_, err = keeper.GetSecret()
	assert.NotNil(t, err)
	val, _ := sqlKeeper.GetSecret()
	assert.True(t, len(val) > 0)
}

func getdb() *sql.DB {
	db, err := sql.Open("mysql", dbAddr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func TestSqlSecretKeeper_GeneratePair(t *testing.T) {
	sqlKeeper.TableName = "abc"
	r := sqlKeeper.GeneratePair()
	assert.Nil(t, r)
	sqlKeeper.TableName = "app"
	r = sqlKeeper.GeneratePair()
	assert.NotNil(t, r)
}
