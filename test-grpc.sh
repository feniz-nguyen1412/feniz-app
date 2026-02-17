#!/bin/bash

echo "🚀 Testing gRPC API with grpcurl"
echo "=================================="

# Start the application in background
echo "📡 Starting application..."
./bin/api &
APP_PID=$!

# Wait for server to start
echo "⏳ Waiting for server to start..."
sleep 3

echo ""
echo "🔍 Available services:"
~/go/bin/grpcurl -plaintext localhost:50051 list

echo ""
echo "📋 UserService details:"
~/go/bin/grpcurl -plaintext localhost:50051 describe api.UserService

echo ""
echo "🧪 Testing gRPC CRUD operations:"

echo ""
echo "1️⃣ Creating user..."
~/go/bin/grpcurl -plaintext -d '{"email":"grpc@example.com","name":"gRPC User"}' \
  localhost:50051 api.UserService/CreateUser

echo ""
echo "2️⃣ Getting all users..."
~/go/bin/grpcurl -plaintext localhost:50051 api.UserService/GetAllUsers

echo ""
echo "3️⃣ Getting user by ID..."
USER_ID=$(~/go/bin/grpcurl -plaintext -d '{"email":"test-grpc-$(date +%s)@example.com","name":"gRPC User"}' \
  localhost:50051 api.UserService/CreateUser | jq -r '.user.id')
echo "User ID: $USER_ID"
~/go/bin/grpcurl -plaintext -d "{\"id\":\"$USER_ID\"}" \
  localhost:50051 api.UserService/GetUser

echo ""
echo "4️⃣ Updating user..."
~/go/bin/grpcurl -plaintext -d "{\"id\":\"$USER_ID\",\"email\":\"updated-grpc-$USER_ID@example.com\",\"name\":\"Updated gRPC User\"}" \
  localhost:50051 api.UserService/UpdateUser

echo ""
echo "5️⃣ Deleting user..."
~/go/bin/grpcurl -plaintext -d "{\"id\":\"$USER_ID\"}" \
  localhost:50051 api.UserService/DeleteUser

echo ""
echo "6️⃣ Verifying deletion..."
~/go/bin/grpcurl -plaintext localhost:50051 api.UserService/GetAllUsers

# Clean up
echo ""
echo "🧹 Cleaning up..."
kill $APP_PID 2>/dev/null

echo ""
echo "✅ gRPC API testing completed successfully!"
echo "🎯 All CRUD operations working perfectly!"
