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

# Determine which MD5 command to use
if command -v md5sum >/dev/null 2>&1; then
    MD5_CMD="md5sum"
elif command -v md5 >/dev/null 2>&1; then
    MD5_CMD="md5 -r"
else
    echo "Neither md5sum nor md5 command found"
    exit 1
fi

# Calculate MD5 sum of all proto files
PROTO_MD5=$(find proto -name "*.proto" -type f -exec $MD5_CMD {} \; | sort | $MD5_CMD | cut -d' ' -f1)
MD5_FILE=".proto.md5"

# Check if MD5 has changed
NEED_GENERATE=1
if [ -f "$MD5_FILE" ]; then
    OLD_MD5=$(cat "$MD5_FILE")
    if [ "$OLD_MD5" == "$PROTO_MD5" ]; then
        NEED_GENERATE=0
        echo "Proto files have not changed, skipping code generation..."
    fi
fi

echo $NEED_GENERATE

# Save current MD5
echo "$PROTO_MD5" > "$MD5_FILE"

# Generate proto files
echo "Generating proto files..."

# Check and install required protoc plugins if needed
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

if ! command -v protoc-gen-openapiv2 &> /dev/null; then
    echo "Installing protoc-gen-openapiv2..."
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
fi

# Generate code
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

go run cmd/generator/main.go -proto="proto/*.proto" -output=out

# Move generated proto files to correct location
mkdir -p out/internal/proto
mv out/app/internal/proto/* out/internal/proto/
rm -rf out/app



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

# Wait for service to finish (it won't unless interrupted)
wait $SERVICE_PID

echo "Done! You can now run the server with: cd out && go run cmd/app/main.go" 