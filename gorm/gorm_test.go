package gorm

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/mysql"
)

func TestOpen(t *testing.T) {
	dsn := "root:root123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4"
	_, err := Open(mysql.Open(dsn))

	assert.Equal(t, err, nil)
}

func TestOpenDsnInvalid(t *testing.T) {
	dsn := "root:root123456@tcp(127.0.0.1:3306)//test?charset=utf8mb4"
	_, err := Open(mysql.Open(dsn))

	assert.NotEqual(t, err, nil)
	assert.MatchRegex(t, err.Error(), `^invalid DSN:`)
}
