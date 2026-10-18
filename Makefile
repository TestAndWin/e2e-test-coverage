BINARY_NAME=e2ecoverage

help:	## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*#/#/' | sed -e 's/#/\n/'

build-ui:	## Build the Vue 3 application for production
	cd ui; npm run format; npm run lint; npm run build

build-swagger:	## Generate the swagger doc
	cd api; swag init -g cmd/coverage/main.go --output docs

update-libs:	## Update Golang and Vue 3 libs (npm with 7 day cooldown, see ui/.ncurc.json)
	cd api; go get -u ./...; go mod tidy
	cd ui; ncu -u --cooldown 7; npm install
	$(MAKE) vulncheck

vulncheck:	## Check the Golang and npm dependencies for known vulnerabilities
	@rc=0; $(MAKE) vulncheck-go || rc=1; $(MAKE) vulncheck-npm || rc=1; exit $$rc

vulncheck-go:	## Check the Golang dependencies for known vulnerabilities
	cd api; go run golang.org/x/vuln/cmd/govulncheck@latest ./...

vulncheck-npm:	## Check the npm dependencies for known vulnerabilities
	cd ui; npm audit

build: build-ui	build-swagger ## Build the Golang binary for Linux and create the Docker image
	cd api; go generate ./ui
	cd api; GOOS=linux GOARCH=amd64 go build -o ../bin/${BINARY_NAME}-linux cmd/coverage/main.go
	docker build -t e2ecoverage .

start:	## Start the docker image in DEV mode
	docker run -e DEV=true --env-file docker_env_vars -p 127.0.0.1:8080:8080 e2ecoverage

docker-tar:	## Create tar.gz from the latest image
	docker save e2ecoverage:latest | gzip > ./bin/e2ecoverage.tar.gz

clean:	## Delete the binary
	go clean
	rm -f bin/${BINARY_NAME}-linux
	rm -f bin/e2ecoverage.tar.gz

start-api: ## Start Golang App directly, useful when working on it
	cd api; DEV=true go run cmd/coverage/main.go

start-ui: ## Start Vue Server in Dev Mode
	cd ui; npm run dev