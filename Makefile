GO_VERSION = 1.17.8
GO := GO111MODULE=on go
GIT_TAG?= $(shell git describe --always --tags)
GOPATH ?= $(shell $(GO) env GOPATH)
GO_NOMOD :=GO111MODULE=off go
GOBIN ?= $(GOPATH)/bin
GOLINT ?= $(GOBIN)/golint
GOSEC ?= $(GOBIN)/gosec

# include .env file and export its env vars
# (-include to ignore error if it does not exist)
-include .env

all: update lint format build

default:
	${MAKE} build
	
# Updates depdendencies
update:
	${GO} mod tidy && cd ./integration_tests && pnpm install && cd ..

# Runs tests
test: 
	${GO} test ./...

# Run ts integration test
integration_test:
	@echo "Integration tests: Running"
	`cd integration_tests && ts-node index.ts`
	@echo "Integration tests: Complete"

# Builds a cross-OS binary
build:
	# Only works for apple silicon
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
	${GO} build \
		-tags netgo -ldflags '-s -w' -o ./bin/server ./cmd/server.go

# Runs a hot reload server using air
run-hot:
	${GOBIN}/air --build.cmd "${GO} build -tags netgo -ldflags '-s -w' -o ./bin/server ./cmd/server.go" --build.bin "./bin/server"

# Runs the binary directly
run-direct:
	./bin/server

# Runs the server directly without hot reloading
run-raw:
	${GO} run ./cmd/server.go

# Clean bin and tmp
clean:
	rm -rf ./bin \
	rm -rf ./tmp

# Runs tests and gets test coverage
check:
	@echo Running Tests and outputting results to ./coverage.out
	${GO} test ./... && go tool cover -html=./coverage.out

# Lints all go files and vets for security issues
lint:
	@echo "LINTING: golint"
	$(GO_NOMOD) get -u golang.org/x/lint/golint
	$(GOLINT) -set_exit_status ./...
	@echo "VETTING"
	$(GO) vet ./...

# Formats all .go files
format:
	@echo "Formatting: gofmt"
	$(GO) fmt ./...

# Builds a docker container
docker-build:
	DOCKER_BUILDKIT=0 docker build \
		-f Dockerfile \
		 -t BLedger:latest .

# Starts a docker contianer
docker-run:
	docker run -it --rm \
		-p 8080:8080 \
		BLedger:latest

# Nukes volumes
docker-nuke:
	docker volume rm -f BLedger_db && docker volume rm -f BLedger_cache

# Nukes images (mainly used when starting from scratch)
docker-prune:
	docker image prune -f