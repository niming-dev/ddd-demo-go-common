package mockdb

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	json "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

// User 用于测试的用户结构
type User struct {
	Id        int
	Name      string
	Height    float64
	Profile   UserProfile `gorm:"type:JSON"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UserProfile struct {
	Photo string
}

func (p *UserProfile) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	*p = UserProfile{}
	err := json.Unmarshal(bytes, p)
	return err
}

func (p UserProfile) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Car 用于测试的汽车结构
type Car struct {
	Id        int
	UserId    int
	Name      string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func TestSQLMockGet(t *testing.T) {
	mock, err := NewSQLMock("testdata/get_user_car_mock.yaml")
	require.NoError(t, err)

	db := mock.GormDB()

	var u User
	expectedUser := User{
		Id:        1,
		Name:      "XiaoMing",
		Height:    1.9,
		Profile:   UserProfile{Photo: "http://www.example.com/logo.png"},
		CreatedAt: cast.ToTime("2022-09-13 15:58:28"),
		UpdatedAt: cast.ToTime("2022-09-13 15:58:28"),
		DeletedAt: time.Time{},
	}
	tx := db.Table("users").Find(&u)
	require.NoError(t, tx.Error)
	require.Equal(t, expectedUser, u)

	var c Car
	expectedCar := Car{
		Id:        1,
		UserId:    1,
		Name:      "id_yxgl",
		Color:     "red",
		CreatedAt: cast.ToTime("2022-09-13 17:36:15"),
		UpdatedAt: cast.ToTime("2022-09-13 09:37:43"),
		DeletedAt: time.Time{},
	}
	tx = db.Table("cars").Find(&c)
	require.NoError(t, tx.Error)
	require.Equal(t, expectedCar, c)
}

func TestSQLMockList(t *testing.T) {
	mock, err := NewSQLMock("testdata/list_users_cars_mock.yaml")
	db := mock.GormDB()

	require.NoError(t, err)
	var us []User
	expectedUsers := []User{
		{
			Id:        1,
			Name:      "XiaoMing",
			Height:    1.9,
			Profile:   UserProfile{Photo: "http://www.example.com/logo.png"},
			CreatedAt: cast.ToTime("2022-09-13 15:58:28"),
			UpdatedAt: cast.ToTime("2022-09-13 15:58:28"),
			DeletedAt: time.Time{},
		},
		{
			Id:        2,
			Name:      "XiaoHong",
			Height:    1.6,
			Profile:   UserProfile{Photo: "http://www.example.com/logo1.png"},
			CreatedAt: cast.ToTime("2022-09-13 19:58:28"),
			UpdatedAt: cast.ToTime("2022-09-13 19:58:28"),
			DeletedAt: time.Time{},
		},
	}
	tx := db.Table("users").Find(&us)
	require.NoError(t, tx.Error)
	require.Equal(t, expectedUsers, us)
}
