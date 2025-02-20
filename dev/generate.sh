#!/bin/bash

# Exit on error
set -e

# Clean output directory
echo "Cleaning output directory..."
rm -rf out/

# Create required directories
echo "Creating directories..."
mkdir -p out/api
mkdir -p out/internal/grpc
mkdir -p out/internal/proto
mkdir -p out/migrations

# Generate proto files
echo "Generating proto files..."
# Install required protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1

# Generate code from proto files
protoc \
    --proto_path=. \
    --proto_path=/usr/local/include \
    --go_out=out \
    --go-grpc_out=out \
    --grpc-gateway_out=out \
    --grpc-gateway_opt=logtostderr=true \
    --grpc-gateway_opt=allow_delete_body=true \
    --grpc-gateway_opt=generate_unbound_methods=true \
    --openapiv2_out=out/api \
    proto/*.proto

# Move generated proto files to correct location
mkdir -p out/internal/proto
mv out/app/internal/proto/* out/internal/proto/
rm -rf out/app

# Generate code
echo "Generating code..."
go run cmd/generator/main.go -proto="proto/*.proto" -output=out

# Format generated code
echo "Formatting generated code..."
find out -name "*.go" -exec go fmt {} \;

# Build and run migrations
cd out

echo "Installing dependencies..."
go mod tidy

# Install goose if not installed
if ! command -v goose &> /dev/null; then
    echo "Installing goose..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

echo "Running migrations..."
# Set database URL for goose
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

# Run migrations
goose -dir migrations up

# Commit and push generated code
echo "Committing and pushing changes..."
cd ..
git add .
git commit -m "Update generated code: $(date '+%Y-%m-%d %H:%M:%S')" || true
git push origin main || true

echo "Starting the service..."
cd out
go run cmd/app/main.go &
SERVICE_PID=$!

# Trap Ctrl+C and kill the service
trap 'echo "Stopping service..."; kill $SERVICE_PID; exit' INT

# Wait for service to finish (it won't unless interrupted)
wait $SERVICE_PID

echo "Done! You can now run the server with: cd out && go run cmd/app/main.go" 