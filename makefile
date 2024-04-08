run:
	@go run cmd/main.go -c configs/config.yaml -e mesh mock

build:
	@GOOS=linux go build -o dist/thingmocker cmd/main.go