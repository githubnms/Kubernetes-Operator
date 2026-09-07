.PHONY: run build test fmt vet

# Run the operator locally (skeleton entrypoint for now)
run:
	go run ./cmd/manager

# Build a binary into bin/manager
build:
	go build -o bin/manager ./cmd/manager

# Run unit tests
test:
	go test ./... -v

# Format code
fmt:
	go fmt ./...

# Static analysis / catch common mistakes
vet:
	go vet ./...