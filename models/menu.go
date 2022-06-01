package models

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"type:varchar(50); not null; default:''"`
	Code        string `gorm:"type:varchar(50); not null; default:''"`
	Parent      uint32 `gorm:"type:bigint(20); not null; default:''"`
	Inheritance string `gorm:"type:varchar(50); not null; default:''"`
	URL         string `gorm:"type:varchar(100); not null; default:''"`
	Component   string `gorm:"type:varchar(50); not null; default:''"`
	Perms       string `gorm:"type:varchar(120); not null; default:''"`
	Type        uint32 `gorm:"type:int(10); not null"`
	Icon        string `gorm:"type:varchar(120); not null; default:''"`
	Order       uint32 `gorm:"type:int(100); not null"`
	Hidden      uint32 `gorm:"type:int(20); not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *Menu) TableName() string {
	return "t_menu"
}
