package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name       string `gorm:"type:varchar(25)"`
	Password   string `gorm:"type:varchar(32); not null; default ''"`
	Salt       string `gorm:"type:char(4); size:4; not null; default ''"`
	Email      string `gorm:"type:varchar(100);unique_index"`
	Profession string `gorm:"type:varchar(255); not null; default ''"`
	Avatar     string `gorm:"type:varchar(255); not null; default ''"`
}

// TableName sets the insert table name for this struct type
func (t *User) TableName() string {
	return "t_user"
}

// TableName sets the insert table name for this struct type
func (t *User) CompareHashAndPassword(password []byte) (bool, error) {
	//Compare the password form and database if match
	bytePassword := []byte(t.Password)
	byteHashedPassword := []byte(password)

	err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	if err != nil {
		return false, err
	}
	return true, nil
}

// TableName sets the insert table name for this struct type
func (t *User) GenerateFromPassword(cost int) error {
	bytePassword := []byte(t.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, cost)
	if err != nil {
		return err
	}
	t.Password = string(hashedPassword)
	return nil
}
