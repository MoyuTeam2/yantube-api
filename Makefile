
.PHONY: build build-proto run clean

PROC_NAME = "yantube-api"

build-proto:
	@echo "Building proto..."
	@protoc --go_out=. --go-grpc_out=. proto/streamserver/stream_server.proto

build: build-proto
	@echo "Building server..."
	@go build -o $(PROC_NAME) .

run: build
	@echo "Running..."
	@./$(PROC_NAME)

clean:
	@echo "Cleaning..."
	@rm -f $(PROC_NAME)
