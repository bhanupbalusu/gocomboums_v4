package model

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID              uint64           `gorm:"primary_key;auto_increment" json:"id"`
	RoleName        string           `gorm:"size:255;not null;unique" json:"role_name"`
	UserRoles       []UserRole       `gorm:"foreignKey:RoleID"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID"`
}

type UserRole struct {
	gorm.Model
	UserID uint64 `gorm:"not null" json:"user_id"`
	RoleID uint64 `gorm:"not null" json:"role_id"`
	User   User   `gorm:"foreignKey:UserID"`
	Role   Role   `gorm:"foreignKey:RoleID"`
}
