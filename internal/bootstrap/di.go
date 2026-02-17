package bootstrap

import (
	"gomono_template/internal/shared/database"
	"gomono_template/internal/shared/kafka"
	"gomono_template/internal/shared/redis"
	"gomono_template/internal/user/application/command"
	"gomono_template/internal/user/application/query"
	"gomono_template/internal/user/domain"
	"gomono_template/internal/user/infrastructure/repository"
	"gomono_template/internal/user/interfaces/http"
	"gomono_template/internal/user/interfaces/grpc"
	"gorm.io/gorm"
)

type Container struct {
	Config             *Config
	DB                 *gorm.DB
	RedisClient        *redis.Client
	KafkaProducer      *kafka.Producer
	KafkaConsumer      *kafka.Consumer
	UserRepo           domain.UserRepository
	CreateUserHandler  *command.CreateUserHandler
	UpdateUserHandler  *command.UpdateUserHandler
	DeleteUserHandler  *command.DeleteUserHandler
	GetUserHandler     *query.GetUserHandler
	GetAllUsersHandler *query.GetAllUsersHandler
	UserHandler        *http.UserHandler
	GRPCUserHandler    *grpc.UserHandler
}

func BuildContainer() (*Container, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.NewPostgres(cfg.DBUrl)
	if err != nil {
		return nil, err
	}

	// Auto-migrate
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers)
	kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topics.UserEvents, cfg.Kafka.GroupID)

	// Create Kafka topics
	err = kafka.CreateTopics(cfg.Kafka.Brokers, cfg.Kafka.Topics.UserEvents)
	if err != nil {
		// Log error but don't fail startup
	}

	userRepo := repository.NewUserRepository(db)

	createUserHandler := command.NewCreateUserHandler(userRepo)
	updateUserHandler := command.NewUpdateUserHandler(userRepo)
	deleteUserHandler := command.NewDeleteUserHandler(userRepo)
	getUserHandler := query.NewGetUserHandler(userRepo)
	getAllUsersHandler := query.NewGetAllUsersHandler(userRepo)

	userHandler := http.NewUserHandler(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
		getUserHandler,
		getAllUsersHandler,
	)

	grpcUserHandler := grpc.NewUserHandler(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
		getUserHandler,
		getAllUsersHandler,
	)

	return &Container{
		Config:             cfg,
		DB:                 db,
		RedisClient:        redisClient,
		KafkaProducer:      kafkaProducer,
		KafkaConsumer:      kafkaConsumer,
		UserRepo:           userRepo,
		CreateUserHandler:  createUserHandler,
		UpdateUserHandler:  updateUserHandler,
		DeleteUserHandler:  deleteUserHandler,
		GetUserHandler:     getUserHandler,
		GetAllUsersHandler: getAllUsersHandler,
		UserHandler:        userHandler,
		GRPCUserHandler:    grpcUserHandler,
	}, nil
}
