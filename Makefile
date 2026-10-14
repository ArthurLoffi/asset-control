.PHONY: swagger test test-integration vet build

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --generalInfo main.go --dir ./cmd/api,./internal --output ./internal/app/swagger --outputTypes go,json,yaml

test:
	go test -race ./...

test-integration:
	go test -race -tags=integration ./...

vet:
	go vet ./...

build:
	go build -o asset-control ./cmd/api
