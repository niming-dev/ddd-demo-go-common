package idgenerator

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrUpdateFailed      = errors.New("update id failed")
	ErrInvalidDatetimeId = errors.New("invalid datetime id")
)

type IdGenerator struct {
	DSN    string // 数据库的连接串
	Prefix string // id前缀，同一个id前缀使用相同的普通Id长度
	Length int    // 普通ID的长度，缺省情况下长度为8
	db     *gorm.DB
}

func CreateIdGenerator(DSN string) (*IdGenerator, error) {
	ret := &IdGenerator{DSN: DSN, Length: 8}
	d, err := gorm.Open(mysql.New(mysql.Config{
		DSN: DSN,
	}), &gorm.Config{})
	if nil != err {
		return nil, err
	}

	sqlDB, err := d.DB()
	if nil != err {
		return nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(20)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	ret.db = d

	d.AutoMigrate(&NormalId{})
	d.AutoMigrate(&DatetimeId{})
	d.AutoMigrate(&DatetimeIdHistory{})

	return ret, nil
}

func (ig *IdGenerator) SetDSN(dsn string) *IdGenerator {
	ig.DSN = dsn
	return ig
}

func (ig *IdGenerator) SetPrefix(prefix string) *IdGenerator {
	ig.Prefix = prefix
	return ig
}

func (ig *IdGenerator) SetLength(length int) *IdGenerator {
	ig.Length = length
	return ig
}

func idPlus(id string, op int, length int) (string, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	if nil != err {
		return "", err
	}
	i += int64(op)
	digitstr := strconv.FormatInt(i, 36)
	digitlen := len(digitstr)
	return strings.Repeat("0", length-digitlen) + digitstr, nil
}

func (ig *IdGenerator) _generateNormalIds(tx *gorm.DB, count int) ([]*NormalId, error) {
	var last NormalId
	var begin string
	err := tx.Model(&NormalId{}).Where("prefix = ?", ig.Prefix).Last(&last).Error
	if nil != err {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			begin = strings.Repeat("0", ig.Length)
		} else {
			return nil, err
		}
	} else {
		begin, err = idPlus(last.IdString, 1, ig.Length)
		if nil != err {
			return nil, err
		}
	}

	var ret []*NormalId
	for i := 0; i < count; i++ {
		ret = append(ret, &NormalId{Prefix: ig.Prefix, IdString: begin, Used: 0})
		begin, err = idPlus(begin, 1, ig.Length)
		if nil != err {
			return nil, err
		}
	}
	err = tx.CreateInBatches(ret, len(ret)).Error
	if nil != err {
		return nil, err
	}

	return ret, nil
}

func (ig *IdGenerator) GetId() (string, error) {
	var id string
	err := ig.db.Transaction(func(tx *gorm.DB) error {
		var normalId *NormalId
		err := tx.Model(&NormalId{}).Where("prefix = ? and used = ?", ig.Prefix, 0).Order("rand()").First(&normalId).Error
		if nil != err {
			if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
				// 产生一批待用id
				ids, err := ig._generateNormalIds(tx, 100)
				if nil != err {
					return err
				}
				normalId = ids[0]
			} else {
				return err
			}
		}

		normalId.Used = 1
		updateResult := tx.Model(&NormalId{}).Where("prefix = ? and id_string = ? and used = ?",
			ig.Prefix, normalId.IdString, 0).Updates(normalId)
		if nil != updateResult.Error {
			return err
		}
		if updateResult.RowsAffected != 1 {
			return ErrUpdateFailed
		}
		id = normalId.IdString
		return nil
	})
	if nil != err {
		return "", err
	}
	return id, nil
}

func (ig *IdGenerator) GetDatetimeId() (string, error) {
	// 先产生一个id
	id_part1 := time.Now().Format("20060102150405")
	id := id_part1 + "_0000"
	dtid := &DatetimeId{
		Prefix:   ig.Prefix,
		IdString: id,
	}
	err := ig.db.Create(dtid).Error
	if nil == err {
		return id, nil
	}

	// 查找最大值并创建新的
	err = ig.db.Transaction(func(tx *gorm.DB) error {
		lastDateTimeId := &DatetimeId{}

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&DatetimeId{}).
			Where("prefix = ? and id_string like ?", ig.Prefix, id_part1+"%").
			Last(&lastDateTimeId).Error
		if nil != err {
			return err
		}
		var part1, part2 string
		n, err := fmt.Sscanf(lastDateTimeId.IdString, "%14s_%4s", &part1, &part2)
		if nil != err {
			return err
		}
		if n != 2 {
			return ErrInvalidDatetimeId
		}
		i, err := strconv.ParseInt(part2, 10, 64)
		if nil != err {
			return err
		}
		i += 1
		newidString := fmt.Sprintf("%s_%.4d", part1, i)

		err = tx.Create(&DatetimeId{Prefix: ig.Prefix, IdString: newidString}).Error
		if nil != err {
			return err
		}
		id = newidString
		return nil
	})
	if nil != err {
		return "", err
	}
	return id, nil
}

func (ig *IdGenerator) GetIdWithPrefix() (string, error) {
	ret, err := ig.GetId()
	if nil != err {
		return "", err
	}
	return ig.Prefix + ret, nil
}

func (ig *IdGenerator) GetDatetimeIdWithPrefix() (string, error) {
	ret, err := ig.GetDatetimeId()
	if nil != err {
		return "", err
	}
	return ig.Prefix + ret, nil
}
