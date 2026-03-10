MODULE_DIR := hw12_13_14_15_16_calendar
BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	cd $(MODULE_DIR) && go build -v -o ../$(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./$(MODULE_DIR)/configs/config.toml

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
	# Создаем минимальный конфиг для golangci-lint 2.x в папке модуля
	cd $(MODULE_DIR) && echo 'version: "2"' > .golangci.yml
	cd $(MODULE_DIR) && echo 'linters: {disable-all: true, enable: [govet]}' >> .golangci.yml
	# Устанавливаем и запускаем линтер
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest
	cd $(MODULE_DIR) && golangci-lint run ./...

.PHONY: build run build-img run-img version test lint