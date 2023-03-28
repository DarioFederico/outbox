BINARY_NAME=outbox

build:
	@go build -o ./bin/${BINARY_NAME} ./cmd/main.go

run:
	@go run ./cmd/main.go

run-release: build
	./bin/{BINARY_NAME}

clean:
	@go clean
	rm ./bin/${BINARY_NAME}

test:
	@go test ./cmd/main.go

swag:
	@swag init -g ./cmd/main.go