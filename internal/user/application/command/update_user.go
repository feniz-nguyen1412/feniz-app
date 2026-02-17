package command

import (
	sharedDomain "gomono_template/internal/shared/domain"
	"gomono_template/internal/user/domain"
	"time"
)

type UpdateUserCommand struct {
	ID    string
	Email string
	Name  string
}

type UpdateUserHandler struct {
	repo domain.UserRepository
}

func NewUpdateUserHandler(repo domain.UserRepository) *UpdateUserHandler {
	return &UpdateUserHandler{repo: repo}
}

func (h *UpdateUserHandler) Handle(cmd UpdateUserCommand) (*domain.User, error) {
	user, err := h.repo.FindByID(cmd.ID)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and if it's already taken
	if cmd.Email != user.Email {
		_, err := h.repo.FindByEmail(cmd.Email)
		if err == nil {
			return nil, sharedDomain.ErrDuplicateEmail
		}
	}

	user.Email = cmd.Email
	user.Name = cmd.Name
	user.UpdatedAt = time.Now()

	err = h.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
