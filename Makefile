BINARY_NAME=e2e-coverage

build:
	# Swagger
	swag init -g cmd/coverage/main.go --output docs
	# Build VUE app and generate ui_gen.go
	go generate ./ui
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY_NAME}-linux cmd/coverage/main.go

build-local:
	# Swagger
	swag init -g cmd/coverage/main.go --output docs
	# Build VUE app and generate ui_gen.go
	go generate ./ui
	# Build for local system
	go build -o ./bin/${BINARY_NAME} cmd/coverage/main.go

run:
	./bin/${BINARY_NAME}

docker:
	docker build -t e2ecoverage .

docker-run:
	docker run --env-file docker_env_vars  -p 127.0.0.1:8080:8080 e2ecoverage
	
build_and_run: build-local run

clean:
	go clean
	rm bin/${BINARY_NAME}
	rm bin/${BINARY_NAME}-win
	rm bin/${BINARY_NAME}-linux