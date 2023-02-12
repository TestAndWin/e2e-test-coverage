BINARY_NAME=e2e-coverage

help:	## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*#/#/' | sed -e 's/#/\n/'

build-ui:	## Build the Vue 3 application for production
	cd ui; npm run build

build-swagger:	## Generate the swagger doc
	cd api; swag init -g cmd/coverage/main.go --output docs

update-libs:	## Update Golang and Vue 3 libs
	cd api; go get -u ./...
	cd ui; ncu -u; npm install

build: build-ui build-swagger	## Build the Golang binary for Linux and create the Docker image
	cd api; go generate ./ui
	cd api; GOOS=linux GOARCH=amd64 go build -o ../bin/${BINARY_NAME}-linux cmd/coverage/main.go
	docker build -t e2ecoverage .

start:	## Start the docker image
	docker run --env-file docker_env_vars  -p 127.0.0.1:8080:8080 e2ecoverage

clean:	## Delete the binary
	go clean
	rm -f bin/${BINARY_NAME}-linux

start-api: ## Start Golang App directly, useful when working on it
	cd api; go run cmd/coverage/main.go

start-ui: ## Start Vue Server in Dev Mode
	cd ui; npm run serve