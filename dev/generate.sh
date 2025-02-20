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
protoc --go_out=out/internal/proto \
    --go_opt=module=app/internal/proto \
    --go-grpc_out=out/internal/proto \
    --go-grpc_opt=module=app/internal/proto \
    proto/*.proto

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