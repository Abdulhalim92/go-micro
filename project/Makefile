FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOG_BINARY=loggerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the broker binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the broker binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOG_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the front end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

PROTO_LOGGER_DIR := ./../logger-service/logs
PROTO_LOGGER_FILE := $(PROTO_LOGGER_DIR)/logs.proto

PROTO_BROKER_DIR := ./../broker-service/logs
PROTO_BROKER_FILE := $(PROTO_BROKER_DIR)/logs.proto

## proto: generate protobuf files
proto_logger:
	@echo "Generating protobuf files..."
	@mkdir -p $(PROTO_LOGGER_DIR)  # Ensure the directory exists
	protoc \
        --proto_path=$(PROTO_LOGGER_DIR) \
        --go_out=$(PROTO_LOGGER_DIR) \
        --go_opt=paths=source_relative \
        --go-grpc_out=$(PROTO_LOGGER_DIR) \
        --go-grpc_opt=paths=source_relative \
        $(PROTO_LOGGER_FILE)



clean_logger:
	@echo "Cleaning generated protobuf files..."
	rm -f $(PROTO_LOGGER_DIR)/*.pb.go
	rm -f $(PROTO_LOGGER_DIR)/*_grpc.pb.go

## proto: generate protobuf files
proto_broker:
	@echo "Generating protobuf files..."
	@mkdir -p $(PROTO_BROKER_DIR)  # Ensure the directory exists
	protoc \
        --proto_path=$(PROTO_BROKER_DIR) \
        --go_out=$(PROTO_BROKER_DIR) \
        --go_opt=paths=source_relative \
        --go-grpc_out=$(PROTO_BROKER_DIR) \
        --go-grpc_opt=paths=source_relative \
        $(PROTO_BROKER_FILE)



clean_broker:
	@echo "Cleaning generated protobuf files..."
	rm -f $(PROTO_BROKER_DIR)/*.pb.go
	rm -f $(PROTO_BROKER_DIR)/*_grpc.pb.go