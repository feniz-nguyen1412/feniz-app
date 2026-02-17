package query

import (
	"gomono_template/internal/user/domain"
)

type GetUserQuery struct {
	ID string
}

type GetUserHandler struct {
	repo domain.UserRepository
}

func NewGetUserHandler(repo domain.UserRepository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

func (h *GetUserHandler) Handle(query GetUserQuery) (*domain.User, error) {
	return h.repo.FindByID(query.ID)
}

type GetAllUsersHandler struct {
	repo domain.UserRepository
}

func NewGetAllUsersHandler(repo domain.UserRepository) *GetAllUsersHandler {
	return &GetAllUsersHandler{repo: repo}
}

func (h *GetAllUsersHandler) Handle() ([]*domain.User, error) {
	return h.repo.FindAll()
}