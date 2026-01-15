package model

type Permission struct {
	ID     int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Permit string `gorm:"column:permit" json:"permit"`
	Name   string `gorm:"column:name" json:"name"`
}

func (*Permission) TableName() string {
	return "permission"
}
