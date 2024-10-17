run:
	@go run cmd/main.go -c configs/config.yaml -e mesh mock

build:
	@GOOS=linux GOARCH=amd64 go build -o dist/thingmocker cmd/main.go

run-local:
	@go run cmd/main.go -c configs/local_config.yaml -e raw mock