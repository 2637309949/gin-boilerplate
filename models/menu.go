package models

import (
	"github.com/jinzhu/gorm"
)

type Menu struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50)"`
	Code        string `gorm:"type:varchar(50)"`
	Parent      uint32 `gorm:"type:bigint(20); not null"`
	Inheritance string `gorm:"type:varchar(50)"`
	URL         string `gorm:"type:varchar(100)"`
	Component   string `gorm:"type:varchar(50)"`
	Perms       string `gorm:"type:varchar(120)"`
	Type        uint32 `gorm:"type:int(10); not null"`
	Icon        string `gorm:"type:varchar(120)"`
	Order       uint32 `gorm:"type:int(100); not null"`
	Hidden      uint32 `gorm:"type:int(20); not null"`
}

func (m *Menu) TableName() string {
	return "t_menu"
}
