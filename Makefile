.PHONY: swagger test test-integration vet build

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --dir ./src/main --output ./src/main/infrastructure/inbound/http/swagger

test:
	go test -race ./...

test-integration:
	go test -race -tags=integration ./...

vet:
	go vet ./...

build:
	go build ./src/main/main.go