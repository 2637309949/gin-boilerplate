package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleMenu struct {
	ID        uint   `gorm:"primarykey"`
	RoleId    uint32 `gorm:"type:bigint(20); not null"`
	MenuId    uint32 `gorm:"type:bigint(20); not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *RoleMenu) TableName() string {
	return "t_role_menu"
}
