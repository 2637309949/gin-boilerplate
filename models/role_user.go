package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleUser struct {
	ID        uint   `gorm:"primarykey"`
	UserId    uint32 `gorm:"type:bigint(20); not null"`
	RoleId    uint32 `gorm:"type:bigint(20); not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName table name of defined RoleUser
func (m *RoleUser) TableName() string {
	return "t_role_user"
}
