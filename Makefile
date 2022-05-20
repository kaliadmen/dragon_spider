## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool dist directory
build:
	@go build -o ./dist/dragon_spider ./cmd/cli

## install_cli: builds the command line tool dragon_spider and installs it to ~/go/bin
install_cli:
	@go build -o ~/go/bin/dragon_spider ./cmd/cli