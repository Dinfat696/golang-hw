MODULE_DIR := hw12_13_14_15_16_calendar
BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	cd $(MODULE_DIR) && go build -v -o ../$(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./$(MODULE_DIR)/configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	cd $(MODULE_DIR) && go test -race ./internal/...

lint:
	cd $(MODULE_DIR) && go mod tidy
	cd $(MODULE_DIR) && go vet ./...
	cd $(MODULE_DIR) && go fmt ./...
	cd $(MODULE_DIR) && go run honnef.co/go/tools/cmd/staticcheck@latest ./...

.PHONY: build run build-img run-img version test lint