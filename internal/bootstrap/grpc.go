package bootstrap

import (
	"gomono_template/internal/user/interfaces/grpc"
	pb "gomono_template/api/proto"
	grpcLib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func NewGRPCServer(container *Container) (*grpcLib.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":"+container.Config.GRPCPort)
	if err != nil {
		panic(err)
	}

	s := grpcLib.NewServer()

	userHandler := grpc.NewUserHandler(
		container.CreateUserHandler,
		container.UpdateUserHandler,
		container.DeleteUserHandler,
		container.GetUserHandler,
		container.GetAllUsersHandler,
	)

	// Register gRPC service
	pb.RegisterUserServiceServer(s, userHandler)

	// Enable reflection for development
	reflection.Register(s)

	return s, lis
}
