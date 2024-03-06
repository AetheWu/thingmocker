run:
	@go run cmd/main.go -c configs/config.yaml -e test mock

build:
	@go build -o dist/thingmocker cmd/main.go