package repository

import (
	sharedDomain "gomono_template/internal/shared/domain"
	"gomono_template/internal/user/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, sharedDomain.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "email = ?", email).Error
	if err == gorm.ErrRecordNotFound {
		return nil, sharedDomain.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}
