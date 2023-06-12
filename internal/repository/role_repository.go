package repository

import (
	"github.com/bhanupbalusu/gocomboums_v4/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(role *model.Role) (*model.Role, error)
	GetRoleByID(id uint64) (*model.Role, error)
	UpdateRole(role *model.Role) (*model.Role, error)
	DeleteRole(id uint64) error
	AddUserRole(userID uint64, roleID uint64) error
	RemoveUserRole(userID uint64, roleID uint64) error
	GetAllRoles() ([]model.Role, error)
	GetRolesByUserID(userID uint64) ([]model.Role, error)
	UserHasRole(userID uint64, roleName string) (bool, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) CreateRole(role *model.Role) (*model.Role, error) {
	if err := r.db.Create(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) GetRoleByID(id uint64) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("UserRoles").Preload("RolePermissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) UpdateRole(role *model.Role) (*model.Role, error) {
	if err := r.db.Save(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) DeleteRole(id uint64) error {
	if err := r.db.Delete(&model.Role{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *roleRepository) AddUserRole(userID uint64, roleID uint64) error {
	userRole := model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	// start transaction
	tx := r.db.Begin()

	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback() // rollback if user role creation fails
		return err
	}

	tx.Commit() // commit transaction if all is well
	return nil
}

func (r *roleRepository) RemoveUserRole(userID uint64, roleID uint64) error {
	// start transaction
	tx := r.db.Begin()

	if err := tx.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&model.UserRole{}).Error; err != nil {
		tx.Rollback() // rollback if deletion fails
		return err
	}

	tx.Commit() // commit transaction if all is well
	return nil
}

func (r *roleRepository) GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.Preload("UserRoles").Preload("RolePermissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) GetRolesByUserID(userID uint64) ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.Joins("JOIN user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) UserHasRole(userID uint64, roleName string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.UserRole{}).Joins("JOIN roles on roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.role_name = ?", userID, roleName).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
