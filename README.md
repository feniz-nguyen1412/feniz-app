
# Go Mono Template

A comprehensive Go monorepo template with clean architecture, CRUD operations, Docker Compose, Redis, Kafka, and gRPC support.

## Project Structure

```
gomono_template/
├── cmd/
│   └── api/
│       └── main.go
│
├── configs/
│   └── config.yaml
│
├── internal/
│   ├── bootstrap/                  # App wiring layer
│   │   ├── config.go
│   │   ├── http.go
│   │   ├── grpc.go
│   │   ├── di.go
│   │   └── module.go
│   │
│   ├── shared/                     # Cross-cutting
│   │   ├── domain/
│   │   │   ├── errors.go
│   │   │   ├── entity.go
│   │   │   └── repository.go
│   │   ├── database/
│   │   │   └── database.go
│   │   ├── redis/
│   │   │   └── redis.go
│   │   └── kafka/
│   │       └── kafka.go
│   │
│   └── user/                       # Bounded Context: User
│       ├── domain/
│       │   ├── user.go
│       │   ├── repository.go
│       │   ├── service.go
│       │   └── errors.go
│       │
│       ├── application/
│       │   ├── command/
│       │   │   ├── create_user.go
│       │   │   ├── update_user.go
│       │   │   └── delete_user.go
│       │   ├── query/
│       │   │   └── get_user.go
│       │   └── dto/
│       │
│       ├── infrastructure/
│       │   └── repository/
│       │       └── user_repository_pg.go
│       │
│       └── interfaces/
│           ├── http/
│           │   ├── handler.go
│           │   └── router.go
│           └── grpc/
│               └── handler.go
│
├── api/
│   └── proto/
│       ├── user.proto
│       ├── user.pb.go
│       └── user_grpc.pb.go
│
├── scripts/
│   └── install-protoc.sh
│
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## Quick Start

### 1. Setup Development Environment

```bash
# Install dependencies and setup tools
make dev-setup

# Start all services (PostgreSQL, Redis, Kafka)
make docker-up

# Build and run the application
make dev
```

### 2. Manual Setup

```bash
# 1. Install dependencies
go mod tidy

# 2. Start Docker services
docker-compose up -d

# 3. Generate protobuf files (if you modify .proto files)
make proto

# 4. Build and run
go build -o bin/api ./cmd/api
./bin/api
```

## Services

### HTTP API Endpoints

- **POST** `http://localhost:8080/users` - Create a new user
- **GET** `http://localhost:8080/users` - Get all users
- **GET** `http://localhost:8080/users/:id` - Get a user by ID
- **PUT** `http://localhost:8080/users/:id` - Update a user
- **DELETE** `http://localhost:8080/users/:id` - Delete a user

### gRPC API

- **Port**: `localhost:50051`
- **Service**: `UserService`
- **Methods**: CreateUser, GetUser, GetAllUsers, UpdateUser, DeleteUser

### Infrastructure Services

- **PostgreSQL**: `localhost:5432`
- **Redis**: `localhost:6379`
- **Kafka**: `localhost:9092`
- **Kafka UI**: `http://localhost:8080` (from kafka-ui container)

## API Examples

### HTTP API Examples

#### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe"
  }'
```

#### Get All Users
```bash
curl http://localhost:8080/users
```

#### Get User by ID
```bash
curl http://localhost:8080/users/{user-id}
```

#### Update User
```bash
curl -X PUT http://localhost:8080/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "email": "updated@example.com",
    "name": "Jane Doe"
  }'
```

#### Delete User
```bash
curl -X DELETE http://localhost:8080/users/{user-id}
```

### gRPC Client Example

**Note**: gRPC service requires proto files to be generated first.

```bash
# 1. Install protoc and plugins
./scripts/install-protoc.sh

# 2. Generate protobuf files
make proto

# 3. Rebuild and run
make build && make run

# 4. Test with grpcurl
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 describe api.UserService

# 5. Test gRPC methods
grpcurl -plaintext -d '{"email":"test@example.com","name":"Test User"}' \
  localhost:50051 api.UserService/CreateUser
```

**Current Status**: gRPC server is running but service is not registered because proto files need to be generated. Run `make proto` to enable full gRPC functionality.

## Docker Services

### PostgreSQL
- **Database**: `gomono_template`
- **User**: `postgres`
- **Password**: `postgres`
- **Port**: `5432`

### Redis
- **Port**: `6379`
- **No authentication** (development only)

### Kafka
- **Brokers**: `localhost:9092`
- **Zookeeper**: `localhost:2181`
- **UI**: `http://localhost:8080` (kafka-ui)
- **Topics**: Auto-created `user_events`

## Configuration

Configuration is managed via `configs/config.yaml`:

```yaml
app:
  port: "8080"
  grpc_port: "50051"

database:
  url: "host=localhost user=postgres password=postgres dbname=gomono_template port=5432 sslmode=disable"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

kafka:
  brokers: ["localhost:9092"]
  group_id: "gomono_template_group"
  topics:
    user_events: "user_events"
```

## Make Commands

```bash
# Build
make build              # Build the application
make run                # Run the application
make clean              # Clean build artifacts

# Protobuf
make proto              # Generate protobuf files

# Docker
make docker-up          # Start Docker services
make docker-down        # Stop Docker services
make logs               # View Docker logs

# Development
make deps               # Install dependencies
make test               # Run tests
make dev-setup          # Full development setup
make dev                # Development workflow
```

## Technology Stack

- **Go 1.18**
- **Gin** - HTTP web framework
- **gRPC** - RPC framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database
- **Redis** - Caching and session storage
- **Kafka** - Event streaming
- **Docker Compose** - Container orchestration
- **Viper** - Configuration management
- **Protocol Buffers** - Service definition

## Architecture

This project follows Clean Architecture principles:

### Layers
- **Domain Layer**: Contains business logic and entities
- **Application Layer**: Contains use cases (commands & queries)
- **Infrastructure Layer**: Contains external dependencies (database, redis, kafka)
- **Interface Layer**: Contains HTTP and gRPC handlers

### Communication Patterns
- **HTTP REST API**: For external client communication
- **gRPC**: For internal service communication
- **Kafka**: For event-driven architecture
- **Redis**: For caching and fast data access

## Features

- ✅ Complete CRUD operations for users
- ✅ Clean architecture with separation of concerns
- ✅ HTTP REST API with Gin
- ✅ gRPC API with Protocol Buffers
- ✅ PostgreSQL integration with GORM
- ✅ Redis integration for caching
- ✅ Kafka integration for event streaming
- ✅ Docker Compose for development environment
- ✅ Dependency injection without external framework
- ✅ Proper error handling and HTTP status codes
- ✅ Email uniqueness validation
- ✅ Configuration management with Viper
- ✅ Auto-migration for database schema
- ✅ Makefile for common operations

## Development Workflow

1. **Start Services**: `make docker-up`
2. **Install Dependencies**: `make deps`
3. **Generate Code**: `make proto` (if .proto files changed)
4. **Run Application**: `make dev`
5. **Test APIs**: Use HTTP or gRPC clients

## Production Considerations

For production deployment:
- Enable Redis authentication
- Use proper Kafka SSL configuration
- Configure PostgreSQL connection pooling
- Add proper logging and monitoring
- Implement rate limiting
- Add API authentication/authorization
- Use environment variables for sensitive config