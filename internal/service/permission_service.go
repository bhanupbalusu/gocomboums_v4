package service

import (
	"strings"

	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"github.com/bhanupbalusu/gocomboums_v4/internal/repository"
	"github.com/bhanupbalusu/gocomboums_v4/pkg/errors"
	"github.com/bhanupbalusu/gocomboums_v4/pkg/logs"
)

type PermissionService struct {
	PermissionRepo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) *PermissionService {
	return &PermissionService{
		PermissionRepo: repo,
	}
}

func validateAndSanitizePermission(permission *model.Permission) error {
	if permission == nil {
		logs.Error("received nil permission", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Permission cannot be nil")
	}

	permission.PermissionName = strings.TrimSpace(strings.ToLower(permission.PermissionName))

	if len(permission.PermissionName) <= 0 {
		logs.Error("permission name is less than three characters", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Permission name cannot be less than three characters")
	}

	if len(permission.PermissionName) > 255 {
		logs.Error("permission name is too long", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Permission name must be less than 256 characters")
	}

	return nil
}

func (s *PermissionService) GetAllPermissions() ([]model.Permission, error) {
	permissions, err := s.PermissionRepo.GetAllPermissions()
	if err != nil {
		logs.Error("error getting all permissions", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return permissions, nil
}

func (s *PermissionService) GetPermissionByID(permissionID uint64) (model.Permission, error) {
	if permissionID <= 0 {
		logs.Error("invalid permission id", nil)
		return model.Permission{}, errors.NewAppError(errors.CodeBadRequest, "Invalid permission id")
	}
	permission, err := s.PermissionRepo.GetPermissionByID(permissionID)
	if err != nil {
		logs.Error("error getting permission by id", err)
		return model.Permission{}, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return permission, nil
}

func (s *PermissionService) CreatePermission(permission model.Permission) (model.Permission, error) {
	// Validate and sanitize permission
	err := validateAndSanitizePermission(&permission)
	if err != nil {
		return model.Permission{}, err
	}

	permission, err = s.PermissionRepo.CreatePermission(permission)
	if err != nil {
		logs.Error("error creating permission", err)
		return model.Permission{}, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return permission, nil
}

func (s *PermissionService) UpdatePermission(permission model.Permission) (model.Permission, error) {
	// Validate and sanitize permission
	err := validateAndSanitizePermission(&permission)
	if err != nil {
		return model.Permission{}, err
	}

	updatedPermission, err := s.PermissionRepo.UpdatePermission(permission)
	if err != nil {
		logs.Error("error updating permission", err)
		return model.Permission{}, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return updatedPermission, nil
}

func (s *PermissionService) DeletePermission(permissionID uint64) error {
	if permissionID <= 0 {
		logs.Error("invalid permission id", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Invalid permission id")
	}
	err := s.PermissionRepo.DeletePermission(permissionID)
	if err != nil {
		logs.Error("error deleting permission", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return nil
}

func (s *PermissionService) AssignPermissionToRole(roleID, permissionID uint64) error {
	if roleID <= 0 || permissionID <= 0 {
		logs.Error("invalid role or permission id", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Invalid role or permission id")
	}
	err := s.PermissionRepo.AssignPermissionToRole(roleID, permissionID)
	if err != nil {
		logs.Error("error assigning permission to role", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return nil
}

func (s *PermissionService) RemovePermissionFromRole(roleID, permissionID uint64) error {
	if roleID <= 0 || permissionID <= 0 {
		logs.Error("invalid role or permission id", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Invalid role or permission id")
	}
	err := s.PermissionRepo.RemovePermissionFromRole(roleID, permissionID)
	if err != nil {
		logs.Error("error removing permission from role", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return nil
}

func (s *PermissionService) GetPermissionsByRoleID(roleID uint64) ([]model.Permission, error) {
	if roleID <= 0 {
		logs.Error("invalid role id", nil)
		return nil, errors.NewAppError(errors.CodeBadRequest, "Invalid role id")
	}
	permissions, err := s.PermissionRepo.GetPermissionsByRoleID(roleID)
	if err != nil {
		logs.Error("error getting permissions by role id", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return permissions, nil
}

func (s *PermissionService) GetRolesByPermissionID(permissionID uint64) ([]model.Role, error) {
	if permissionID <= 0 {
		logs.Error("invalid permission id", nil)
		return nil, errors.NewAppError(errors.CodeBadRequest, "Invalid permission id")
	}
	roles, err := s.PermissionRepo.GetRolesByPermissionID(permissionID)
	if err != nil {
		logs.Error("error getting roles by permission id", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return roles, nil
}

func (s *PermissionService) AddMultiplePermissionsToRole(roleID uint64, permissionIDs []uint64) error {
	if roleID <= 0 || len(permissionIDs) == 0 {
		logs.Error("invalid role id or empty permissions", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Invalid role id or empty permissions")
	}
	err := s.PermissionRepo.AddMultiplePermissionsToRole(roleID, permissionIDs)
	if err != nil {
		logs.Error("error adding multiple permissions to role", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return nil
}

func (s *PermissionService) RemoveMultiplePermissionsFromRole(roleID uint64, permissionIDs []uint64) error {
	if roleID <= 0 || len(permissionIDs) == 0 {
		logs.Error("invalid role id or empty permissions", nil)
		return errors.NewAppError(errors.CodeBadRequest, "Invalid role id or empty permissions")
	}
	err := s.PermissionRepo.RemoveMultiplePermissionsFromRole(roleID, permissionIDs)
	if err != nil {
		logs.Error("error removing multiple permissions from role", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred.")
	}
	return nil
}
