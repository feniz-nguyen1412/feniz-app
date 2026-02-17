package grpc

import (
	"context"
	sharedDomain "gomono_template/internal/shared/domain"
	"gomono_template/internal/user/application/command"
	"gomono_template/internal/user/application/query"
	"gomono_template/internal/user/domain"
	pb "gomono_template/api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	createHandler  *command.CreateUserHandler
	updateHandler  *command.UpdateUserHandler
	deleteHandler  *command.DeleteUserHandler
	getHandler     *query.GetUserHandler
	getAllHandler *query.GetAllUsersHandler
}

func NewUserHandler(
	createHandler *command.CreateUserHandler,
	updateHandler *command.UpdateUserHandler,
	deleteHandler *command.DeleteUserHandler,
	getHandler *query.GetUserHandler,
	getAllHandler *query.GetAllUsersHandler,
) *UserHandler {
	return &UserHandler{
		createHandler:  createHandler,
		updateHandler:  updateHandler,
		deleteHandler:  deleteHandler,
		getHandler:     getHandler,
		getAllHandler:  getAllHandler,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := h.createHandler.Handle(command.CreateUserCommand{
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		if err == sharedDomain.ErrDuplicateEmail {
			return nil, status.Error(codes.AlreadyExists, "Email already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUserResponse{
		User: domainUserToProto(user),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.getHandler.Handle(query.GetUserQuery{ID: req.Id})
	if err != nil {
		if err == sharedDomain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserResponse{
		User: domainUserToProto(user),
	}, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := h.getAllHandler.Handle()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoUsers []*pb.User
	for _, user := range users {
		protoUsers = append(protoUsers, domainUserToProto(user))
	}

	return &pb.GetAllUsersResponse{
		Users: protoUsers,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err := h.updateHandler.Handle(command.UpdateUserCommand{
		ID:    req.Id,
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		if err == sharedDomain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		if err == sharedDomain.ErrDuplicateEmail {
			return nil, status.Error(codes.AlreadyExists, "Email already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateUserResponse{
		User: domainUserToProto(user),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.deleteHandler.Handle(command.DeleteUserCommand{ID: req.Id})
	if err != nil {
		if err == sharedDomain.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteUserResponse{}, nil
}

func domainUserToProto(user *domain.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
