#!/bin/bash

# Exit on error
set -e

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "protoc is not installed. Please install it first."
    exit 1
fi

# Install required Go tools if not installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install grpc-gateway
cd out && go get github.com/grpc-ecosystem/grpc-gateway/v2 && cd ..

# Clean output directory
echo "Cleaning output directory..."
rm -rf out/
mkdir -p out/

# Create required directories
echo "Creating directories..."
mkdir -p out/cmd/app
mkdir -p out/internal/handler
mkdir -p out/internal/service
mkdir -p out/internal/repository
mkdir -p out/internal/models
mkdir -p out/internal/grpc
mkdir -p out/internal/tests
mkdir -p out/migrations

# Generate proto files
echo "Generating proto files..."
# Generate gRPC code from proto files
protoc -I . \
    -I$(go list -f '{{.Dir}}' -m github.com/grpc-ecosystem/grpc-gateway/v2)/runtime/internal/examplepb \
    -I$(go list -f '{{.Dir}}' -m github.com/grpc-ecosystem/grpc-gateway/v2)/.. \
    --go_out=out \
    --go_opt=paths=source_relative \
    --go-grpc_out=out \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=out \
    --grpc-gateway_opt=paths=source_relative \
    proto/service.proto

# Move generated proto files to correct location
mkdir -p out/internal/proto
mv out/proto/* out/internal/proto/
rm -rf out/proto

# Generate application code
go run cmd/generator/main.go -proto proto/service.proto -output out/

# Format generated code
echo "Formatting generated code..."
cd out && go fmt ./... && cd ..

# Build and run migrations
cd out

echo "Installing dependencies..."
go mod tidy

# Run migrations
echo "Running migrations..."
# Wait for PostgreSQL to be ready
until PGPASSWORD=postgres psql -h localhost -U postgres -c '\q' 2>/dev/null; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

# Run migrations down
echo "Rolling back all migrations..."
GOOSE_DRIVER=postgres \
  GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" \
  goose -dir migrations down-to 0

# Run migrations up
echo "Applying all migrations..."
GOOSE_DRIVER=postgres \
  GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" \
  goose -dir migrations up

echo "Starting the service..."
go run cmd/app/main.go &
SERVICE_PID=$!

# Trap Ctrl+C and kill the service
trap 'echo "Stopping service..."; kill $SERVICE_PID; exit' INT

# Wait for service to finish (it won't unless interrupted)
wait $SERVICE_PID

echo "Done! You can now run the server with: cd out && go run cmd/app/main.go" 