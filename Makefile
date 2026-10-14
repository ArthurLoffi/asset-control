.PHONY: swagger test test-integration vet build-linux build-win

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --generalInfo main.go --dir ./cmd/api,./internal --output ./internal/app/swagger --outputTypes go,json,yaml

test:
	go test -race ./...

test-integration:
	go test -race -tags=integration ./...

vet:
	go vet ./...

# Build for linux
build-linux:
	go build -o asset-control ./cmd/api

#Build for windows
build-win:
	go build -o asset-control.exe ./cmd/api