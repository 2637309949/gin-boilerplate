package models

import (
	"github.com/jinzhu/gorm"
)

type RolePermission struct {
	gorm.Model
	RoleId       uint32 `gorm:"type:bigint(20); not null"`
	PermissionId uint32 `gorm:"type:bigint(20); not null"`
}

func (m *RolePermission) TableName() string {
	return "t_role_permission"
}
