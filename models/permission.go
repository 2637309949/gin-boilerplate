package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(25); not null; default:''"`
	Code      string `gorm:"type:varchar(25); not null; default:''"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Permission) TableName() string {
	return "t_permission"
}
