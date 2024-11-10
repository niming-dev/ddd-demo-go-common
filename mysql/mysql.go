package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	ufgorm "github.com/niming-dev/ddd-demo/go-common/gorm"
)

func Open(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	dialector := mysql.Open(dsn)
	return ufgorm.Open(dialector, opts...)
}
