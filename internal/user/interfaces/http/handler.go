package http

import (
	sharedDomain "gomono_template/internal/shared/domain"
	"gomono_template/internal/user/application/command"
	"gomono_template/internal/user/application/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createHandler  *command.CreateUserHandler
	updateHandler  *command.UpdateUserHandler
	deleteHandler  *command.DeleteUserHandler
	getHandler     *query.GetUserHandler
	getAllHandler  *query.GetAllUsersHandler
}

func NewUserHandler(
	createHandler *command.CreateUserHandler,
	updateHandler *command.UpdateUserHandler,
	deleteHandler *command.DeleteUserHandler,
	getHandler *query.GetUserHandler,
	getAllHandler *query.GetAllUsersHandler,
) *UserHandler {
	return &UserHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		getAllHandler: getAllHandler,
	}
}

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

type UpdateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.createHandler.Handle(command.CreateUserCommand{
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		if err == sharedDomain.ErrDuplicateEmail {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	
	user, err := h.getHandler.Handle(query.GetUserQuery{ID: id})
	if err != nil {
		if err == sharedDomain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.getAllHandler.Handle()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.updateHandler.Handle(command.UpdateUserCommand{
		ID:    id,
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		if err == sharedDomain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err == sharedDomain.ErrDuplicateEmail {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	
	err := h.deleteHandler.Handle(command.DeleteUserCommand{ID: id})
	if err != nil {
		if err == sharedDomain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}