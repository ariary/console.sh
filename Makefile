GOFILES=`go list ./...`
GOFILESNOTEST=`go list ./... | grep -v test`

before.build:
	go mod download && go mod vendor

build.console.sh: lint
	@echo "build in ${PWD}";go build cmd/console.sh/console.sh.go

build.console.sh.static: lint
	@echo "build in ${PWD}";CGO_ENABLED=0 go build cmd/console.sh/console.sh.go

lint: ## Lint the files
	@go fmt ${GOFILES}
	@go vet ${GOFILESNOTEST}