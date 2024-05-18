
.PHONY: build build-proto run clean deps test build-go

PROC_NAME = "yantube-api"

build-proto:
	@echo "Building proto..."
	@protoc --go_out=. --go-grpc_out=. proto/streamserver/stream_server.proto

build-go:
	@echo "Building server..."
	@go build -o $(PROC_NAME) .

build: build-proto build-go

deps:
	@echo "Installing dependencies..."
	@go mod download

run: build-go
	@echo "Running..."
	@./$(PROC_NAME)

clean:
	@echo "Cleaning..."
	@rm -f $(PROC_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...