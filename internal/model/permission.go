package model

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	ID              uint64           `gorm:"primary_key;auto_increment" json:"id"`
	PermissionName  string           `gorm:"size:255;not null;unique" json:"permission_name"`
	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID"`
}

type RolePermission struct {
	gorm.Model
	RoleID       uint64     `gorm:"not null" json:"role_id"`
	PermissionID uint64     `gorm:"not null" json:"permission_id"`
	Role         Role       `gorm:"foreignKey:RoleID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
}
