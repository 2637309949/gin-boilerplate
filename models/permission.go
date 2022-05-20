package models

import (
	"github.com/jinzhu/gorm"
)

type Permission struct {
	gorm.Model
	Name string `gorm:"type:varchar(25)"`
	Code string `gorm:"type:varchar(25)"`
}

func (m *Permission) TableName() string {
	return "t_permission"
}
