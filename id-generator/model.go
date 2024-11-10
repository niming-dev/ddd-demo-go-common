package idgenerator

type NormalId struct {
	Id int `gorm:"primaryKey;autoIncrement"`
	// 前缀，某个前缀下所有id长度应该一致
	Prefix string `gorm:"size:16;uniqueIndex:idx_normal_id_prefix_string,priority:1"`
	// 后台产生一堆数据插入
	IdString string `gorm:"size:32;uniqueIndex:idx_normal_id_prefix_string,priority:2"`
	// 0表示未使用，其他表示已使用
	Used int `gorm:"default:0;index"`
}

func (NormalId) TableName() string {
	return "t_normal_id"
}

// 此表中有数据的说明Id已经被用过了
// 表中的老数据会移动到历史表，IdString必须采用20201122010203_221这种格式，避免跟历史的冲突
type DatetimeId struct {
	Id int `gorm:"primaryKey;autoIncrement"`
	// 前缀，某个前缀下id唯一
	Prefix   string `gorm:"size:16;uniqueIndex:idx_datetime_id_prefix_string,priority:1"`
	IdString string `gorm:"size:32;uniqueIndex:idx_datetime_id_prefix_string,priority:2"`
}

func (DatetimeId) TableName() string {
	return "t_datetime_id"
}

type DatetimeIdHistory struct {
	Id int `gorm:"primaryKey;autoIncrement"`
	// 前缀，某个前缀下id唯一
	Prefix   string `gorm:"size:16;uniqueIndex:idx_datetime_id_hi_prefix_string,priority:1"`
	IdString string `gorm:"size:32;uniqueIndex:idx_datetime_id_hi_prefix_string,priority:2"`
}

func (DatetimeIdHistory) TableName() string {
	return "t_datetime_id_hi"
}
