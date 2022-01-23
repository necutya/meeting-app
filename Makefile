export GO111MODULE=on
export GOSUMDB=off

BIN_NAME := $(or $(PROJECT_NAME), 'meeting-app')
GOLINT := golangci-lint

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1

dep: # Download required dependencies
	go mod tidy
	go mod download
	go mod vendor

cilint: dep check-lint
	$(GOLINT) run -c .golangci.yml --timeout 5m

lint: cilint

build: dep
	CGO_ENABLED=1 go build -mod=vendor -o ./bin/${BIN_NAME} -a ./cmd/meeting-app

run: dep
	go run ./cmd/meeting-app

clean: ## Remove previous build
	rm -f bin/$(BIN_NAME)

check-swagger:
	@which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

external-swagger: check-swagger
	swagger generate spec -o src/server/http/static/ext/swagger.yaml  -w ./ --scan-models --include-tag=external
	swagger generate spec -o src/server/http/static/ext/swagger.json  -w ./ --scan-models --include-tag=external

swagger: check-swagger external-swagger
	swagger generate spec -o src/server/http/static/v1/swagger.yaml  -w ./ --scan-models --exclude-tag=external
	swagger generate spec -o src/server/http/static/v1/swagger.json  -w ./ --scan-models --exclude-tag=external

fmt: ## format source files
	go fmt gitlab.yalantis.com/erp/api/src/...
