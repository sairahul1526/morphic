.PHONY: format clean vet lint build server test docs setup check docker-up docker-down migrate

format:
	go fmt ./...

clean:
	rm -f ./morphic
	GO111MODULE=on go mod tidy -v

vet:
	go vet ./...

lint:
	golangci-lint run --config=.golangci.yml

build:
	go build -o morphic main.go

server:
	go run main.go server start

migrate:
	go run main.go migrate

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docs:
	swag init
	promlinter list metrics/metric.go --output md --add-help > docs/metrics.md

setup:
	go install github.com/yeya24/promlinter/cmd/promlinter@v0.3.0

check: setup clean format vet lint docs test

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

all: check docker-up migrate