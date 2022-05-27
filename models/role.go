package models

import (
	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Name       string `gorm:"type:varchar(25); not null; default:''"`
	Code       string `gorm:"type:varchar(25); not null; default:''"`
	Status     uint32 `gorm:"type:int(10); not null; default:0"`
	AppIndex   string `gorm:"type:varchar(50); not null; default:''"`
	AdminIndex string `gorm:"type:varchar(50); not null; default:''"`
}

// TableName table name of defined Role
func (m *Role) TableName() string {
	return "t_role"
}
