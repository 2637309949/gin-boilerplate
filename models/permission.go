package models

import (
	"github.com/jinzhu/gorm"
)

type Permission struct {
	gorm.Model
	Name string `gorm:"type:varchar(25); not null; default:''"`
	Code string `gorm:"type:varchar(25); not null; default:''"`
}

func (m *Permission) TableName() string {
	return "t_permission"
}
