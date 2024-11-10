package mysql

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestOpen(t *testing.T) {
	dsn := "root:root123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4"
	_, err := Open(dsn)

	assert.Equal(t, err, nil)
}

func TestOpenDsnInvalid(t *testing.T) {
	dsn := "root:root123456@tcp(127.0.0.1:3306)//test?charset=utf8mb4"
	_, err := Open(dsn)

	assert.NotEqual(t, err, nil)
	assert.MatchRegex(t, err.Error(), `^invalid DSN:`)
}
