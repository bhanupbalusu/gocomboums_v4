package repository

import (
	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"gorm.io/gorm"
)

// PermissionRepository interface
type PermissionRepository interface {
	GetAllPermissions() ([]model.Permission, error)
	GetPermissionByID(permissionID uint64) (model.Permission, error)
	CreatePermission(permission model.Permission) (model.Permission, error)
	UpdatePermission(permission model.Permission) (model.Permission, error)
	DeletePermission(permissionID uint64) error
	AssignPermissionToRole(roleID, permissionID uint64) error
	RemovePermissionFromRole(roleID, permissionID uint64) error
	GetPermissionsByRoleID(roleID uint64) ([]model.Permission, error)
	GetRolesByPermissionID(permissionID uint64) ([]model.Role, error)
	AddMultiplePermissionsToRole(roleID uint64, permissionIDs []uint64) error
	RemoveMultiplePermissionsFromRole(roleID uint64, permissionIDs []uint64) error
}

// permissionRepository struct
type permissionRepository struct {
	DBConn *gorm.DB
}

// NewPermissionRepository returns a new permissionRepository instance
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		DBConn: db,
	}
}

// GetAllPermissions gets all permissions from the database
func (repo *permissionRepository) GetAllPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	if err := repo.DBConn.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermissionByID gets a permission by its ID
func (repo *permissionRepository) GetPermissionByID(permissionID uint64) (model.Permission, error) {
	var permission model.Permission
	if err := repo.DBConn.First(&permission, permissionID).Error; err != nil {
		return model.Permission{}, err
	}
	return permission, nil
}

// CreatePermission creates a new permission
func (repo *permissionRepository) CreatePermission(permission model.Permission) (model.Permission, error) {
	if err := repo.DBConn.Create(&permission).Error; err != nil {
		return model.Permission{}, err
	}
	return permission, nil
}

// UpdatePermission updates a permission
func (repo *permissionRepository) UpdatePermission(permission model.Permission) (model.Permission, error) {
	if err := repo.DBConn.Save(&permission).Error; err != nil {
		return model.Permission{}, err
	}
	return permission, nil
}

// DeletePermission deletes a permission by its ID
func (repo *permissionRepository) DeletePermission(permissionID uint64) error {
	if err := repo.DBConn.Delete(&model.Permission{}, permissionID).Error; err != nil {
		return err
	}
	return nil
}

// AssignPermissionToRole assigns a permission to a role
func (repo *permissionRepository) AssignPermissionToRole(roleID, permissionID uint64) error {
	tx := repo.DBConn.Begin()

	rolePermission := model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	if err := tx.Create(&rolePermission).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// RemovePermissionFromRole removes a permission from a role
func (repo *permissionRepository) RemovePermissionFromRole(roleID, permissionID uint64) error {
	tx := repo.DBConn.Begin()

	rolePermission := model.RolePermission{}

	if err := tx.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&rolePermission).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&rolePermission).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetPermissionsByRoleID gets permissions by role ID
func (repo *permissionRepository) GetPermissionsByRoleID(roleID uint64) ([]model.Permission, error) {
	rolePermissions := []model.RolePermission{}
	if err := repo.DBConn.Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	permissions := []model.Permission{}
	for _, rolePermission := range rolePermissions {
		permission := model.Permission{}
		if err := repo.DBConn.First(&permission, rolePermission.PermissionID).Error; err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (repo *permissionRepository) GetRolesByPermissionID(permissionID uint64) ([]model.Role, error) {
	rolePermissions := []model.RolePermission{}
	if err := repo.DBConn.Where("permission_id = ?", permissionID).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	roles := []model.Role{}
	for _, rolePermission := range rolePermissions {
		role := model.Role{}
		if err := repo.DBConn.First(&role, rolePermission.RoleID).Error; err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (repo *permissionRepository) AddMultiplePermissionsToRole(roleID uint64, permissionIDs []uint64) error {
	tx := repo.DBConn.Begin()

	for _, permissionID := range permissionIDs {
		rolePermission := model.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}

		if err := tx.Create(&rolePermission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (repo *permissionRepository) RemoveMultiplePermissionsFromRole(roleID uint64, permissionIDs []uint64) error {
	tx := repo.DBConn.Begin()

	for _, permissionID := range permissionIDs {
		rolePermission := model.RolePermission{}

		if err := tx.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&rolePermission).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(&rolePermission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
