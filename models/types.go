package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint `gorm:"not null;auto_increment=false;primary_key;unique_index:composite_index"`
	Username string
	Email    string
}

type File struct {
	gorm.Model
	ID       uint
	FileName string
	Data     []byte
	OwnerID  uint
}
