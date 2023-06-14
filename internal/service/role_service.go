package service

import (
	"strings"

	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"github.com/bhanupbalusu/gocomboums_v4/internal/repository"
	"github.com/bhanupbalusu/gocomboums_v4/pkg/errors"
	"github.com/bhanupbalusu/gocomboums_v4/pkg/logs"
)

type RoleService struct {
	RoleRepo repository.RoleRepository
}

// NewUserService creates a new UserService with the provided repo
func NewRoleService(repo repository.RoleRepository) *RoleService {
	return &RoleService{
		RoleRepo: repo,
	}
}

func validateAndSanitizeRole(role *model.Role) error {
	// Validate input
	if role == nil {
		return errors.NewAppError(errors.CodeBadRequest, "role cannot be nil")
	}

	if len(role.RoleName) < 3 || len(role.RoleName) > 255 {
		logs.Error("role name length is out of allowed range", nil)
		return errors.NewAppError(errors.CodeBadRequest, "role name length must be between 3 and 255 characters")
	}

	// Sanitize input
	role.RoleName = strings.TrimSpace(role.RoleName)

	return nil
}

func (s *RoleService) CreateRole(role *model.Role) (*model.Role, error) {
	// Validate and Sanitize input
	err := validateAndSanitizeRole(role)
	if err != nil {
		return nil, err
	}

	newRole, err := s.RoleRepo.CreateRole(role)
	if err != nil {
		logs.Error("error creating role", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return newRole, nil
}

func (s *RoleService) GetRoleByID(id uint64) (*model.Role, error) {
	// Check if id is valid
	if id == 0 {
		logs.Error("invalid role id", errors.NewAppError(errors.CodeBadRequest, "invalid role id"))
		return nil, errors.NewAppError(errors.CodeBadRequest, "invalid role id")
	}

	role, err := s.RoleRepo.GetRoleByID(id)
	if err != nil {
		logs.Error("error fetching role by id", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return role, nil
}

func (s *RoleService) UpdateRole(role *model.Role) (*model.Role, error) {
	err := validateAndSanitizeRole(role)
	if err != nil {
		return nil, err
	}

	updatedRole, err := s.RoleRepo.UpdateRole(role)
	if err != nil {
		logs.Error("error updating role", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return updatedRole, nil
}

func (s *RoleService) DeleteRole(id uint64) error {
	// Check if id is valid
	if id == 0 {
		logs.Error("invalid role id", errors.NewAppError(errors.CodeBadRequest, "invalid role id"))
		return errors.NewAppError(errors.CodeBadRequest, "Username already exists")
	}

	err := s.RoleRepo.DeleteRole(id)
	if err != nil {
		logs.Error("error deleting role", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return nil
}

func (s *RoleService) AddUserRole(userID uint64, roleID uint64) error {
	// Validate user id and role id
	if userID == 0 {
		logs.Error("invalid user id", errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero"))
		return errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero")
	}

	if roleID == 0 {
		logs.Error("invalid role id", errors.NewAppError(errors.CodeBadRequest, "role id cannot be zero"))
		return errors.NewAppError(errors.CodeBadRequest, "role id cannot be zero")
	}

	err := s.RoleRepo.AddUserRole(userID, roleID)
	if err != nil {
		logs.Error("error adding role to user", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return nil
}

func (s *RoleService) RemoveUserRole(userID uint64, roleID uint64) error {
	// Validate user id and role id
	if userID == 0 {
		logs.Error("invalid user id", errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero"))
		return errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero")
	}

	if roleID == 0 {
		logs.Error("invalid role id", errors.NewAppError(errors.CodeBadRequest, "role id cannot be zero"))
		return errors.NewAppError(errors.CodeBadRequest, "role id cannot be zero")
	}

	err := s.RoleRepo.RemoveUserRole(userID, roleID)
	if err != nil {
		logs.Error("error removing role from user", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return nil
}

func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	roles, err := s.RoleRepo.GetAllRoles()
	if err != nil {
		logs.Error("error fetching all roles", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return roles, nil
}

func (s *RoleService) GetRolesByUserID(userID uint64) ([]model.Role, error) {
	// Validate user id
	if userID == 0 {
		logs.Error("invalid user id", errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero"))
		return nil, errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero")
	}

	roles, err := s.RoleRepo.GetRolesByUserID(userID)
	if err != nil {
		logs.Error("error fetching roles by user id", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return roles, nil
}

func (s *RoleService) UserHasRole(userID uint64, roleName string) (bool, error) {
	// Validate user id and role name
	if userID == 0 {
		logs.Error("invalid user id", errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero"))
		return false, errors.NewAppError(errors.CodeBadRequest, "user id cannot be zero")
	}

	if len(roleName) < 2 {
		logs.Error("role name too short", errors.NewAppError(errors.CodeBadRequest, "role name must be at least 2 characters long"))
		return false, errors.NewAppError(errors.CodeBadRequest, "role name must be at least 2 characters long")
	}

	hasRole, err := s.RoleRepo.UserHasRole(userID, roleName)
	if err != nil {
		logs.Error("error checking user role", err)
		return false, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return hasRole, nil
}
