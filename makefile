run:
	@go run ./... -c configs/config.yaml -e development mock

build:
	@GOOS=linux go build -o thingmocker ./...