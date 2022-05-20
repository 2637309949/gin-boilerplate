package models

import (
	"github.com/jinzhu/gorm"
)

type RoleUser struct {
	gorm.Model
	UserId uint32 `gorm:"type:bigint(20); not null"`
	RoleId uint32 `gorm:"type:bigint(20); not null"`
}

// TableName table name of defined RoleUser
func (m *RoleUser) TableName() string {
	return "t_role_user"
}
