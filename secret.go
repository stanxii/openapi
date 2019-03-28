package openapi

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// the interface to get the secret
type SecretKeeper interface {
	GetSecret() (string, error)
}

const (
	queryTemplate  = "SELECT `%s` FROM `%s` WHERE `%s` = ? LIMIT 1"
	insertTemplate = "INSERT INTO `%s`(`%s`, `%s`)VALUE(?, ?)"
	EmptyString    = ""
)

/**
# the sql to create the table app
CREATE TABLE `app` (
  `app_key` varchar(32) NOT NULL,
  `app_secret` varchar(128) NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
 */

// default provided sql
type SqlSecretKeeper struct {
	Db        *sql.DB // the client to access database
	TableName string  // the table where the secret stores
	KeyCol    string  // the column name of the key
	AppKey    string  // the app key of the api
	SecretCol string  // the column name of the secret
}

// get secret from a sql data source
func (s SqlSecretKeeper) GetSecret() (string, error) {
	if s.Db == nil {
		return EmptyString, errors.New("db client should not be nil")
	}
	row := s.Db.QueryRow(s.constructQuery(), s.AppKey)
	var secret string
	err := row.Scan(&secret)
	if err != nil {
		return EmptyString, err
	}
	return secret, nil
}

// construct query for getting secret
func (s SqlSecretKeeper) constructQuery() string {
	return fmt.Sprintf(queryTemplate, s.SecretCol, s.TableName, s.KeyCol)
}

func (s SqlSecretKeeper) GeneratePair() *KvPair {
	p := KvPair{
		Key:   string(randomStr(keyLen, kindAll)),
		Value: string(randomStr(secretLen, kindAll)),
	}
	// do the insert work
	insertSql := fmt.Sprintf(insertTemplate, s.TableName, s.KeyCol, s.SecretCol)
	r, err := s.Db.Exec(insertSql, p.Key, p.Value)
	if err != nil || r == nil {
		return nil
	}
	a, err := r.RowsAffected()
	// check result
	if err != nil || a < 1 {
		return nil
	}
	return &p
}

// generate random string
const (
	keyLen    = 16
	secretLen = 32
	kindAll   = 3
)

// export RandStr
// random strings
func randomStr(size int, kind int) []byte {
	kinds := [][]int{{10, 48}, {26, 97}, {26, 65}}
	specialChars := []byte{95, 45, 46, 35, 36, 37, 38}
	specialCharLen := len(specialChars)
	iKind, result := kind, make([]byte, size)
	isAll := kind == 3
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {

		// random iKind
		if isAll {
			iKind = rand.Intn(3)
		}
		if iKind == 3 {
			result[i] = specialChars[rand.Intn(specialCharLen)]
		} else {
			scope, base := kinds[iKind][0], kinds[iKind][1]
			result[i] = uint8(base + rand.Intn(scope))
		}
	}
	return result
}
