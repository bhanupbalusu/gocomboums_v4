package service

import (
	"fmt"
	"strings"

	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"github.com/bhanupbalusu/gocomboums_v4/internal/repository"
	"github.com/bhanupbalusu/gocomboums_v4/pkg/errors"
	"github.com/bhanupbalusu/gocomboums_v4/pkg/logs"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo repository.UserRepository
}

// NewUserService creates a new UserService with the provided repo
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func validateInput(user *model.User) error {
	if user == nil {
		logs.Error("Received nil user for creation")
		return errors.NewAppError(errors.CodeBadRequest, "User cannot be nil")
	}

	if len(user.Username) < 3 || len(user.Username) > 255 {
		logs.Error("Username length is out of allowed range")
		return errors.NewAppError(errors.CodeBadRequest, "Username length must be between 3 and 255 characters")
	}

	if len(user.PasswordHash) < 6 {
		logs.Error("Password too short")
		return errors.NewAppError(errors.CodeBadRequest, "Password must be at least 6 characters long")
	}

	if !strings.Contains(user.Email, "@") {
		logs.Error("Invalid email address")
		return errors.NewAppError(errors.CodeBadRequest, "Email must be a valid address")
	}

	return nil
}

func sanitizeInput(user *model.User) {
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
}

func HashPassword(user *model.User) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		logs.Error("Error hashing password", err)
		return nil, err
	}

	return hashedPassword, err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *UserService) CreateUser(user *model.User) error {
	// Validate input
	if err := validateInput(user); err != nil {
		return errors.NewAppError(errors.CodeBadRequest, "Invalid user data")
	}

	// Sanitize input
	sanitizeInput(user)

	// Check if username already exists
	if existingUser, err := s.UserRepo.GetUserByUsername(user.Username); err != nil {
		logs.Error("Error fetching user by username", err)
		return errors.NewAppError(errors.CodeInternalServerError, "An unexpected error occurred")
	} else if existingUser != nil {
		return errors.NewAppError(errors.CodeBadRequest, "Username already exists")
	}

	// Create the user
	if err := s.UserRepo.CreateUser(user); err != nil {
		logs.Error("Error creating user", err)
		return errors.NewAppError(errors.CodeInternalServerError, "An unexpected error occurred")
	}

	return nil
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		logs.Error("error fetching user by username: ", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}
	return user, nil
}

func (s *UserService) GetUserByID(id uint64) (*model.User, error) {
	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		logs.Error("error fetching user by id: ", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}
	return user, nil
}

func (s *UserService) UpdateUser(user *model.User) error {
	// Validate input
	if user.ID == 0 {
		logs.Error("User ID cannot be zero for update")
		return errors.NewAppError(errors.CodeBadRequest, "User ID cannot be zero")
	}

	// Validate input
	err := validateInput(user)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	// Sanitize input
	sanitizeInput(user)

	// Check if username already exists
	existingUser, err := s.UserRepo.GetUserByUsername(user.Username)
	if err != nil {
		logs.Error("Error fetching user by username", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}
	if existingUser != nil && existingUser.ID != user.ID {
		logs.Error(fmt.Sprintf("Username already exists: %s", user.Username))
		return errors.NewAppError(errors.CodeBadRequest, "Username already exists")
	}

	// Update the user
	err = s.UserRepo.UpdateUser(user)
	if err != nil {
		logs.Error("Error updating user", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return nil
}

func (s *UserService) DeleteUser(id uint64) error {
	// Validate input
	if id == 0 {
		logs.Error("User ID cannot be zero for delete")
		return errors.NewAppError(errors.CodeBadRequest, "User ID cannot be zero")
	}

	// Check if user exists
	existingUser, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		logs.Error("Error fetching user by id", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}
	if existingUser == nil {
		logs.Error(fmt.Sprintf("User does not exist for id: %d", id))
		return errors.NewAppError(errors.CodeNotFound, "User not found")
	}

	// Delete the user
	err = s.UserRepo.DeleteUser(id)
	if err != nil {
		logs.Error("Error deleting user", err)
		return errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return nil
}

func (s *UserService) ListUsers(page int, pageSize int) ([]*model.User, error) {
	if page < 0 || pageSize <= 0 {
		logs.Error("Invalid pagination parameters", errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred"))
		return nil, errors.NewAppError(errors.CodeBadRequest, "User ID cannot be zero")
	}

	users, err := s.UserRepo.ListUsers(page, pageSize)
	if err != nil {
		logs.Error("Failed to fetch users", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return users, nil
}

func (s *UserService) SearchUsers(query string, page int, pageSize int) ([]*model.User, error) {
	if page < 0 || pageSize <= 0 {
		logs.Error("Invalid pagination parameters", errors.NewAppError(errors.CodeBadRequest, "User ID cannot be zero"))
		return nil, errors.NewAppError(errors.CodeBadRequest, "User ID cannot be zero")
	}

	users, err := s.UserRepo.SearchUsers(query, page, pageSize)
	if err != nil {
		logs.Error("Failed to search users", err)
		return nil, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return users, nil
}

func (s *UserService) CountUsers() (int64, error) {
	count, err := s.UserRepo.CountUsers()
	if err != nil {
		logs.Error("Failed to count users", err)
		return 0, errors.NewAppError(errors.CodeInternalServerError, "Internal server error occurred")
	}

	return count, nil
}
