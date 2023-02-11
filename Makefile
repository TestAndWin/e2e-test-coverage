BINARY_NAME=e2e-coverage

build-ui:
	cd ui; npm run build

build-swagger:
	cd api; swag init -g cmd/coverage/main.go --output docs

update-libs:
	cd api; go get -u ./...
	cd ui; ncu -u; npm install

build: build-ui build-swagger
	# Build VUE app and generate ui_gen.go
	cd api; go generate ./ui
	cd api; GOOS=linux GOARCH=amd64 go build -o ../bin/${BINARY_NAME}-linux cmd/coverage/main.go

docker:
	docker build -t e2ecoverage .

docker-run:
	docker run --env-file docker_env_vars  -p 127.0.0.1:8080:8080 e2ecoverage

build-local: build-ui build-swagger
	# Build VUE app and generate ui_gen.go
	cd api; go generate ./ui
	# Build for local system
	cd api; go build -o ../bin/${BINARY_NAME} cmd/coverage/main.go

run:
	./bin/${BINARY_NAME}

build-and-run: build-local run

clean:
	go clean
	rm -f bin/${BINARY_NAME}
	rm -f bin/${BINARY_NAME}-linux

# Start Golang App directly
run-local-api:
	cd api; go run cmd/coverage/main.go

# Start VUE Server in Dev Mode
run-local-ui:
	cd ui; npm run serve