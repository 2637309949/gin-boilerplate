package models

import (
	"time"

	"gorm.io/gorm"
)

type RolePermission struct {
	ID           uint   `gorm:"primarykey"`
	RoleId       uint32 `gorm:"type:bigint(20); not null"`
	PermissionId uint32 `gorm:"type:bigint(20); not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (m *RolePermission) TableName() string {
	return "t_role_permission"
}
