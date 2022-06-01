package models

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"primarykey"`
	Name       string `gorm:"type:varchar(25)"`
	Password   string `gorm:"type:varchar(32); not null; default:''"`
	Salt       string `gorm:"type:char(4); size:4; not null; default:''"`
	Email      string `gorm:"type:varchar(100);uniqueIndex"`
	Profession string `gorm:"type:varchar(255); not null; default:''"`
	Avatar     string `gorm:"type:varchar(255); not null; default:''"`
	Verified   uint32 `gorm:"type:int(10); not null; default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName sets the insert table name for this struct type
func (t *User) TableName() string {
	return "t_user"
}

// TableName sets the insert table name for this struct type
func (t *User) CompareHashAndPassword(password string) bool {
	//Compare the password form and database if match
	hashPassword, err := t.GenerateFromPassword(password)
	if err != nil {
		return false
	}
	// match
	fmt.Println(hashPassword)
	return hashPassword == t.Password
}

// TableName sets the insert table name for this struct type
func (t *User) GenerateFromPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, []byte(t.Salt)...)
	_, err := sha512Hasher.Write(passwordBytes)
	if err != nil {
		return "", err
	}
	// Convert the hashed password to a base64 encoded string
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash, nil
}
