package entities

type Branch struct {
	Id   uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name;type:varchar(100);not null" json:"name"`
}

func (Branch) TableName() string {
	return "branches"
}

type Department struct {
	Id   uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name;type:varchar(100);not null" json:"name"`
}

func (Department) TableName() string {
	return "departments"
}
