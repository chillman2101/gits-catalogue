package repository

import (
	"github.com/chillman2101/gits-catalogue/internal/model"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByIDAndTokenHash(id uint, tokenHash string) (*model.User, error)
	FindByIDAndRefreshTokenHash(id uint, refreshTokenHash string) (*model.User, error)
	Create(user *model.User) error
	UpdateTokenHash(id uint, hash string) error
	UpdateRefreshTokenHash(id uint, hash string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, res_err.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, res_err.ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepository) FindByIDAndTokenHash(id uint, tokenHash string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND token_hash = ?", id, tokenHash).First(&user).Error
	if err != nil {
		return nil, res_err.ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepository) FindByIDAndRefreshTokenHash(id uint, refreshTokenHash string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND refresh_token_hash = ?", id, refreshTokenHash).First(&user).Error
	if err != nil {
		return nil, res_err.ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return res_err.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdateTokenHash(id uint, hash string) error {
	err := r.db.Model(&model.User{}).Where("id = ?", id).Update("token_hash", hash).Error
	if err != nil {
		return res_err.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdateRefreshTokenHash(id uint, hash string) error {
	err := r.db.Model(&model.User{}).Where("id = ?", id).Update("refresh_token_hash", hash).Error
	if err != nil {
		return res_err.ErrUserNotFound
	}
	return nil
}
