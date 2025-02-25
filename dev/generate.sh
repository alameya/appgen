#!/bin/bash

# Exit on error
set -e

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "protoc is not installed. Please install it first."
    exit 1
fi

# Install required Go tools if not installed
# Проверяем и устанавливаем необходимые инструменты Go
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
fi

if ! command -v protoc-gen-grpc-gateway &> /dev/null; then
    echo "Installing protoc-gen-grpc-gateway..."
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
fi

if ! command -v goose &> /dev/null; then
    echo "Installing goose..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

# Install grpc-gateway
if [ -d "out" ]; then
    cd out && go get github.com/grpc-ecosystem/grpc-gateway/v2 && cd ..
fi

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
mkdir -p out/internal/proto

# Generate proto files
echo "Generating proto files..."
# Проверяем, что файлы существуют
echo "Proto files found:"
ls -la proto/*.proto

# Generate gRPC code from proto files
PROTO_FILES=$(find proto -name '*.proto')
if [ -z "$PROTO_FILES" ]; then
  echo "Error: No proto files found in proto/"
  exit 1
fi

echo "Processing proto files:"
echo "$PROTO_FILES"

protoc -I . \
    -I$(go list -f '{{.Dir}}' -m github.com/grpc-ecosystem/grpc-gateway/v2)/runtime/internal/examplepb \
    -I$(go list -f '{{.Dir}}' -m github.com/grpc-ecosystem/grpc-gateway/v2)/.. \
    --go_out=out/internal \
    --go_opt=paths=source_relative \
    --go-grpc_out=out/internal \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=out/internal \
    --grpc-gateway_opt=paths=source_relative \
    $(echo "$PROTO_FILES" | tr '\n' ' ')

# Move generated proto files to correct location


# Generate application code
# Находим все proto файлы и передаем их генератору
PROTO_FILES=$(find proto -name '*.proto' | tr '\n' ',')
echo "Found proto files: $PROTO_FILES"
go run cmd/generator/main.go -proto "$PROTO_FILES" -output out/

# Format generated code
echo "Formatting generated code..."
cd out && go fmt ./... && cd ..

# Build and run migrations
cd out

# Load environment variables from .env
echo "Loading environment variables..."
if [ -f ../.env ]; then
  export $(cat ../.env | grep -v '^#' | xargs)
elif [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
else
  echo "Error: .env file not found"
  exit 1
fi

echo "Installing dependencies..."
go mod tidy

# Run migrations
echo "Running migrations..."
# Wait for PostgreSQL to be ready
until PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USER} -c '\q' 2>/dev/null; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

# Run migrations down
echo "Rolling back all migrations..."
PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -c "
  DROP SCHEMA public CASCADE;
  CREATE SCHEMA public;
  GRANT ALL ON SCHEMA public TO ${DB_USER};
  GRANT ALL ON SCHEMA public TO public;
"

# Run migrations up
echo "Applying all migrations..."
GOOSE_DRIVER=postgres \
  GOOSE_DBSTRING="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
  goose -dir migrations up

echo "Starting the service..."
go run cmd/app/main.go &
SERVICE_PID=$!

# Trap Ctrl+C and kill the service
trap 'echo "Stopping service..."; kill $SERVICE_PID; exit' INT

# Wait for service to finish (it won't unless interrupted)
wait $SERVICE_PID

echo "Done! You can now run the server with: cd out && go run cmd/app/main.go" 