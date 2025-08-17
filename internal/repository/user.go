package repository

import (
	"live-chat-backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}
