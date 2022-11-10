export GO111MODULE=on

lint:
	golangci-lint run ./...

run:
	go run cmd/coshkey-server/main.go

build:
	go mod download && CGO_ENABLED=0  go build \
		-o ./bin/coshkey-server$(shell go env GOEXE) ./cmd/coshkey-server/main.go

# запуск контейнеров
docker-init: docker-pull docker-up

docker-pull:
	docker pull golang:1.19-alpine
	docker pull alpine:latest

docker-up:
	docker compose up -d