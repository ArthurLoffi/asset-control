.PHONY: swagger

swagger:
	swag init --dir ./src/main --output ./src/main/adapter/inbound/http/swagger
