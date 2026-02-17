#!/bin/bash

# Install protoc and plugins for Go gRPC

echo "Installing protoc..."
if command -v apt-get &> /dev/null; then
    # Ubuntu/Debian
    sudo apt-get update
    sudo apt-get install -y protobuf-compiler
elif command -v yum &> /dev/null; then
    # CentOS/RHEL
    sudo yum install -y protobuf-compiler
elif command -v brew &> /dev/null; then
    # macOS
    brew install protobuf
else
    echo "Please install protoc manually for your OS"
    exit 1
fi

echo "Installing Go protobuf plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Adding $HOME/go/bin to PATH if not already there..."
if ! echo $PATH | grep -q "$HOME/go/bin"; then
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
    export PATH=$PATH:$HOME/go/bin
fi

echo "Installation complete!"
echo "Run 'source ~/.bashrc' or restart your terminal to use the new PATH"
