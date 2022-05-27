package models

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Title         string `gorm:"type:varchar(100); not null; default: '' "`
	Introduction  string `gorm:"type:varchar(255); not null; default: ''"`
	ContentMd     string `gorm:"type:text"`
	ContentHtml   string `gorm:"type:text"`
	DirectoryHtml string `gorm:"type:text"`
	UserID        int    `gorm:"type:int(10); not null; default: 0"`
	User          *User
	Tags          string `gorm:"type:varchar(255); not null; default: '' "`
	ViewNum       int64  `gorm:"type:int(10); not null; default: 0"`
}

// TableName sets the insert table name for this struct type
func (t *Article) TableName() string {
	return "t_article"
}
