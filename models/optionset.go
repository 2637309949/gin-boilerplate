package models

import "github.com/jinzhu/gorm"

type Optionset struct {
	gorm.Model
	Name  string `gorm:"type:varchar(100); not null; default ''"`
	Code  string `gorm:"type:varchar(100); not null; default ''"`
	Value string `gorm:"type:varchar(520); not null; default ''"`
}

// TableName sets the insert table name for this struct type
func (t *Optionset) TableName() string {
	return "t_optionset"
}
