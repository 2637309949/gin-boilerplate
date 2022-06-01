package models

import (
	"time"

	"gorm.io/gorm"
)

type Optionset struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(100); not null; default: ''"`
	Code      string `gorm:"type:varchar(100); not null; default: ''"`
	Value     string `gorm:"type:varchar(520); not null; default: ''"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName sets the insert table name for this struct type
func (t *Optionset) TableName() string {
	return "t_optionset"
}
