package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Username     string     `gorm:"size:255;not null;unique" json:"username"`
	PasswordHash string     `gorm:"size:255;not null;" json:"password_hash"`
	Email        string     `gorm:"size:255;not null;unique" json:"email"`
	UserRoles    []UserRole `gorm:"foreignKey:UserID"`
}
