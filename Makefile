BINARY_NAME=e2e-coverage

build:
	# Swagger
	swag init -g cmd/coverage/main.go --output docs
	# Build VUE app and generate ui_gen.go
	go generate ./ui
	# Build for local system
	go build -o ./bin/${BINARY_NAME} cmd/coverage/main.go
	# Windows build
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY_NAME}-win cmd/coverage/main.go
	# Linux build
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY_NAME}-linux cmd/coverage/main.go

run:
	./bin/${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm bin/${BINARY_NAME}
	rm bin/${BINARY_NAME}-win
	rm bin/${BINARY_NAME}-linux