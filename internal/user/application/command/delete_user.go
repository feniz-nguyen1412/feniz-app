package command

import (
	"gomono_template/internal/user/domain"
)

type DeleteUserCommand struct {
	ID string
}

type DeleteUserHandler struct {
	repo domain.UserRepository
}

func NewDeleteUserHandler(repo domain.UserRepository) *DeleteUserHandler {
	return &DeleteUserHandler{repo: repo}
}

func (h *DeleteUserHandler) Handle(cmd DeleteUserCommand) error {
	// Check if user exists
	_, err := h.repo.FindByID(cmd.ID)
	if err != nil {
		return err
	}

	return h.repo.Delete(cmd.ID)
}
