package models

import (
	"github.com/jinzhu/gorm"
)

type RoleMenu struct {
	gorm.Model
	RoleId uint32 `gorm:"type:bigint(20); not null"`
	MenuId uint32 `gorm:"type:bigint(20); not null"`
}

func (m *RoleMenu) TableName() string {
	return "t_role_menu"
}
