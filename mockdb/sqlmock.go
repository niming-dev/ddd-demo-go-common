package mockdb

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	json "github.com/json-iterator/go"
	"github.com/niming-dev/ddd-demo/go-common/mockdb/reader"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SQLMock .
type SQLMock struct {
	Tables  map[string]SQLMockTable `yaml:"tables" validate:"dive,required"`
	Queries []SQLMockQuery          `yaml:"queries" validate:"dive,required"`

	sqlmock.Sqlmock
	db *sql.DB
	// 配置文件所在路径
	confPath string
}

func NewSQLMock(filename string) (*SQLMock, error) {
	// 文件所在路径
	filePath := path.Dir(filename)
	// 读取文件内容
	bts, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解析模板数据
	tpl, err := NewTemplate(filePath, ".yaml")
	if err != nil {
		return nil, err
	}
	tpl, err = tpl.Parse(string(bts))
	if err != nil {
		return nil, err
	}
	btsBuff := bytes.NewBuffer(nil)
	err = tpl.Execute(btsBuff, nil)
	if err != nil {
		return nil, err
	}

	var mock SQLMock
	err = yaml.Unmarshal(btsBuff.Bytes(), &mock)
	if err != nil {
		return nil, err
	}

	// 保存配置文件所在路径
	mock.confPath = filePath

	// 检查验证规则
	if err = validator.New().Struct(&mock); err != nil {
		return nil, err
	}

	// 初始化mock
	mock.db, mock.Sqlmock, err = sqlmock.New()
	if err != nil {
		return nil, err
	}
	// 关闭SQL执行顺序的检查功能
	mock.MatchExpectationsInOrder(false)

	// 把配置数据注入到mock中
	err = mock.inject()
	if err != nil {
		return nil, err
	}

	return &mock, nil
}

func (s SQLMock) inject() error {
	for _, query := range s.Queries {
		table, ok := s.Tables[query.TableName]
		if !ok {
			return errors.New("undefined Table " + query.TableName)
		}
		rows, err := query.GetRows(s.confPath, table)
		if err != nil {
			return err
		}
		s.ExpectQuery(query.ExpectQuerySQL).WillReturnRows(rows)
	}
	return nil
}

func (s SQLMock) DB() *sql.DB {
	return s.db
}

func (s SQLMock) GormDB() *gorm.DB {
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      s.db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	return gormDB
}

// SQLMockTable 表定义，字段列表
type SQLMockTable struct {
	Table   map[string]*SQLMockColumn `validate:"dive,required"`
	Columns []*SQLMockColumn          `validate:"dive,required"`
}

func (t *SQLMockTable) UnmarshalYAML(node *yaml.Node) error {
	err := node.Decode(&t.Columns)
	if err != nil {
		return err
	}

	t.Table = make(map[string]*SQLMockColumn, len(t.Columns))
	for _, c := range t.Columns {
		t.Table[c.Name] = c
	}
	return nil
}

// checkRow 检查row中的数据是否符合table的定义
func (t *SQLMockTable) checkRow(row map[string]driver.Value) error {
	for k, v := range row {
		column, ok := t.Table[k]
		if !ok {
			return errors.New("undefined column " + k)
		}

		cv, err := column.Check(v)
		if err != nil {
			return err
		}
		row[k] = cv
	}
	return nil
}

// checkRows 检查多个row数据
func (t *SQLMockTable) checkRows(rows []map[string]driver.Value) error {
	for _, row := range rows {
		if err := t.checkRow(row); err != nil {
			return err
		}
	}
	return nil
}

// ColumnNames 返回字段列表
func (t *SQLMockTable) ColumnNames() []string {
	var columns []string
	for _, c := range t.Columns {
		columns = append(columns, c.Name)
	}
	return columns
}

// rowToArray 根据columns列表，把值转换为数组类型
func (t *SQLMockTable) rowToArray(row map[string]driver.Value) (list []driver.Value) {
	for _, c := range t.ColumnNames() {
		list = append(list, row[c])
	}
	return
}

func (t *SQLMockTable) NewEmptyRows() *sqlmock.Rows {
	columns := make([]*sqlmock.Column, len(t.Columns))

	for i, column := range t.Columns {
		c := sqlmock.NewColumn(column.Name).
			Nullable(column.Nullable).
			OfType(string(column.DBType), column.DBType.sameValue())
		if column.Length > 0 {
			c.WithLength(column.Length)
		}
		if column.Precision > 0 && column.Scale > 0 {
			c.WithPrecisionAndScale(column.Precision, column.Scale)
		}
		columns[i] = c
	}
	return sqlmock.NewRowsWithColumnDefinition(columns...)
}

// SQLMockColumn 字段定义
type SQLMockColumn struct {
	Name      string `yaml:"name" validate:"required"`    // 字段名称
	DBType    DBType `yaml:"db_type" validate:"required"` // 类型
	Nullable  bool   `yaml:"nullable"`                    // 是否可为 null
	Length    int64  `yaml:"length"`                      // 长度
	Precision int64  `yaml:"precision"`                   // 用于 decimal 类型，精度
	Scale     int64  `yaml:"scale"`                       // 用于 decimal 类型，小数位数
}

func (c SQLMockColumn) Check(value any) (driver.Value, error) {
	if value == nil {
		return nil, nil
	}
	switch c.DBType {
	case DBTypeChar, DBTypeVarchar, DBTypeTinytext, DBTypeText, DBTypeMediumtext, DBTypeLongtext:
		return cast.ToStringE(value)
	case DBTypeTinyint, DBTypeSmallint, DBTypeMediumint, DBTypeInt, DBTypeInteger, DBTypeBigint:
		return cast.ToInt64E(value)
	case DBTypeBool, DBTypeBoolean:
		return cast.ToBoolE(value)
	case DBTypeFloat, DBTypeDouble, DBTypeDecimal:
		return cast.ToFloat64E(value)
	case DBTypeDate, DBTypeDatetime, DBTypeTimestamp, DBTypeTime, DBTypeYear:
		return cast.ToTimeE(value)
	case DBTypeJSON:
		return json.Marshal(value)
	}

	return nil, fmt.Errorf("invalid type: %t", value)
}

type DBType string

const (
	DBTypeChar       DBType = "CHAR"
	DBTypeVarchar    DBType = "VARCHAR"
	DBTypeTinytext   DBType = "TINYTEXT"
	DBTypeText       DBType = "TEXT"
	DBTypeMediumtext DBType = "MEDIUMTEXT"
	DBTypeLongtext   DBType = "LONGTEXT"
	DBTypeBool       DBType = "BOOL"
	DBTypeBoolean    DBType = "BOOLEAN"
	DBTypeTinyint    DBType = "TINYINT"
	DBTypeSmallint   DBType = "SMALLINT"
	DBTypeMediumint  DBType = "MEDIUMINT"
	DBTypeInt        DBType = "INT"
	DBTypeInteger    DBType = "INTEGER"
	DBTypeBigint     DBType = "BIGINT"
	DBTypeFloat      DBType = "FLOAT"
	DBTypeDouble     DBType = "DOUBLE"
	DBTypeDecimal    DBType = "DECIMAL"
	DBTypeDate       DBType = "DATE"
	DBTypeDatetime   DBType = "DATETIME"
	DBTypeTimestamp  DBType = "TIMESTAMP"
	DBTypeTime       DBType = "TIME"
	DBTypeYear       DBType = "YEAR"
	DBTypeJSON       DBType = "JSON"
)

func (d DBType) sameValue() any {
	switch d {
	case DBTypeChar, DBTypeVarchar, DBTypeTinytext, DBTypeText, DBTypeMediumtext, DBTypeLongtext:
		return "example"
	case DBTypeTinyint, DBTypeSmallint, DBTypeMediumint, DBTypeInt, DBTypeInteger, DBTypeBigint:
		return int64(1)
	case DBTypeBool, DBTypeBoolean:
		return true
	case DBTypeFloat, DBTypeDouble, DBTypeDecimal:
		return 1.1
	case DBTypeDate, DBTypeDatetime, DBTypeTimestamp, DBTypeTime, DBTypeYear:
		return time.Now()
	case DBTypeJSON:
		return `{"json": "example"}`
	}

	return nil
}

// SQLMockQuery 查询定义
type SQLMockQuery struct {
	TableName      string   `yaml:"table_name" validate:"required"`       // 表名
	ExpectQuerySQL string   `yaml:"expect_query_sql" validate:"required"` // 预期执行的SQL，支持正则
	DataFiles      []string `yaml:"data_files" validate:"required"`       // 数据文件集，文件支持的格式：csv
}

func (q SQLMockQuery) GetRows(filePath string, table SQLMockTable) (*sqlmock.Rows, error) {
	rows := table.NewEmptyRows()

	for _, dataFile := range q.DataFiles {
		filename := path.Join(filePath, dataFile)
		dataValues, err := q.readRowsFromFile(filename)
		if err != nil {
			return nil, err
		}

		// 检查数据类型
		err = table.checkRows(dataValues)
		if err != nil {
			return nil, err
		}

		for _, value := range dataValues {
			rows.AddRow(table.rowToArray(value)...)
		}
	}

	return rows, nil
}

type Reader interface {
	Read() ([]map[string]driver.Value, error)
}

// readRowsFromFile 从文件中读取rows数据
func (q SQLMockQuery) readRowsFromFile(filename string) ([]map[string]driver.Value, error) {
	r, err := q.newReader(filename)
	if err != nil {
		return nil, err
	}

	rows, err := r.Read()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (q SQLMockQuery) newReader(filename string) (Reader, error) {
	suffix := strings.ToLower(path.Ext(filename))
	switch suffix {
	case ".json":
		return reader.NewJsonReader(filename), nil
	case ".yaml", ".yml":
		return reader.NewYamlReader(filename), nil
	}

	return nil, errors.New("unsupported suffix: " + suffix)
}
