package domain

import (
	"gomono_template/internal/shared/domain"
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id, email, name string) (*User, error) {
	if email == "" {
		return nil, domain.ErrInvalidEmail
	}

	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
