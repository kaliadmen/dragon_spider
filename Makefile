## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool dragon_spider and copies it to myapp
build_cli:
	@go build -o ../myApp/dragon_spider ./cmd/cli

## build_cli: builds the command line tool dist directory
build:
	@go build -o ./dist/dragon_spider ./cmd/cli