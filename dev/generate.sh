#!/bin/bash

# Exit on error
set -e

# Clean output directory
echo "Cleaning output directory..."
rm -rf out/

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

echo "Running migrations..."
# Используем существующую базу postgres
PGPASSWORD=postgres psql -h localhost -U postgres -d postgres -f migrations/*.sql

# Run migrations
./scripts/migrate.sh up

echo "Starting the service..."
go run cmd/main.go &
SERVICE_PID=$!

# Trap Ctrl+C and kill the service
trap 'echo "Stopping service..."; kill $SERVICE_PID; exit' INT

# Wait for service to finish (it won't unless interrupted)
wait $SERVICE_PID

echo "Done! You can now run the server with: cd out && go run cmd/main.go" 