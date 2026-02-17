package command

import (
	sharedDomain "gomono_template/internal/shared/domain"
	"gomono_template/internal/user/domain"
	"github.com/google/uuid"
)

type CreateUserCommand struct {
	Email string
	Name  string
}

type CreateUserHandler struct {
	repo domain.UserRepository
}

func NewCreateUserHandler(repo domain.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (h *CreateUserHandler) Handle(cmd CreateUserCommand) (*domain.User, error) {
	// Check if user already exists
	_, err := h.repo.FindByEmail(cmd.Email)
	if err == nil {
		return nil, sharedDomain.ErrDuplicateEmail
	}

	user, err := domain.NewUser(uuid.New().String(), cmd.Email, cmd.Name)
	if err != nil {
		return nil, err
	}

	err = h.repo.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
