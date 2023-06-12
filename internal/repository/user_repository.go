// internal/repository/user_repository.go

package repository

import (
	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id uint64) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint64) error
	ListUsers(page int, pageSize int) ([]*model.User, error)
	SearchUsers(query string, page int, pageSize int) ([]*model.User, error)
	CountUsers() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint64) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) ListUsers(page int, pageSize int) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) SearchUsers(query string, page int, pageSize int) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Where("username LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
