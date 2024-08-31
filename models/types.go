package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint
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
