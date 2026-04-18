APP_NAME := task-queue-platform
BUILD_DIR := bin

.PHONY: build run test fmt line clean

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./... -v -race -count=1

fmt:
	gofmt -w
	go vet ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)
	go clean