BINARY_NAME=e2ecoverage
# Keep in sync with github.com/swaggo/swag in api/go.mod
SWAG_VERSION=v1.16.6
# Pinned instead of @latest; only bump to releases older than COOLDOWN_DAYS
GOVULNCHECK_VERSION=v1.9.0
# Supply chain cooldown: only use library versions published at least this many days ago.
# npm side is also enforced by ui/.npmrc (min-release-age) and ui/.ncurc.json (cooldown).
COOLDOWN_DAYS=7
# api/ui/ui.go embeds dist/* relative to itself
EMBED_DIR=api/ui/dist

help:	## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*#/#/' | sed -e 's/#/\n/'

build-ui:	## Build the Vue 3 application for production
	cd ui; npm run format; npm run lint; npm run build

build-swagger:	## Generate the swagger doc
	cd api; go run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION) init -g cmd/coverage/main.go --output docs

update-libs:	## Update Golang and Vue 3 libs, only to versions older than COOLDOWN_DAYS
	cd api; go run ../scripts/gocooldown.go -days $(COOLDOWN_DAYS)
	cd ui; ncu -u --cooldown $(COOLDOWN_DAYS); npm install
	$(MAKE) cooldown-check
	$(MAKE) vulncheck

cooldown-check:	## Check that all locked Go and npm dependencies are older than COOLDOWN_DAYS
	@rc=0; \
	(cd api; go run ../scripts/gocooldown.go -check -days $(COOLDOWN_DAYS)) || rc=1; \
	(cd ui; node ../scripts/npmcooldown.mjs $(COOLDOWN_DAYS)) || rc=1; \
	exit $$rc

vulncheck:	## Check the Golang and npm dependencies for known vulnerabilities
	@rc=0; $(MAKE) vulncheck-go || rc=1; $(MAKE) vulncheck-npm || rc=1; exit $$rc

vulncheck-go:	## Check the Golang dependencies for known vulnerabilities
	cd api; go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

vulncheck-npm:	## Check the npm dependencies for known vulnerabilities
	cd ui; npm audit

embed-ui: build-ui	## Copy the built UI to where api/ui/ui.go expects it for go:embed
	rm -rf ${EMBED_DIR}
	mkdir -p ${EMBED_DIR}
	cp -R ui/dist/. ${EMBED_DIR}/

build: embed-ui	build-swagger ## Build the Golang binary for Linux and create the Docker image
	cd api; GOOS=linux GOARCH=amd64 go build -o ../bin/${BINARY_NAME}-linux cmd/coverage/main.go
	docker build -t e2ecoverage .

start:	## Start the docker image in DEV mode
	docker run -e DEV=true --env-file docker_env_vars -p 127.0.0.1:8080:8080 e2ecoverage

docker-tar:	## Create tar.gz from the latest image
	docker save e2ecoverage:latest | gzip > ./bin/e2ecoverage.tar.gz

clean:	## Delete the binary and all generated UI assets
	go clean
	rm -f bin/${BINARY_NAME}-linux
	rm -f bin/e2ecoverage.tar.gz
	rm -rf ${EMBED_DIR} ui/dist

start-api: embed-ui	## Start Golang App directly, useful when working on it
	cd api; DEV=true go run cmd/coverage/main.go

start-ui: ## Start Vue Server in Dev Mode
	cd ui; npm run dev