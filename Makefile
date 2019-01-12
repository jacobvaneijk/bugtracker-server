GO_PKGS = $(shell go list ./... | sed 's/github.com\/jacobvaneijk\/bugtracker-server/./')
GO_FILES = $(shell find . -name "*.go" | grep -v "_test")
BIN_NAME = "bugtracker-server"

.PHONY: run
run:
	@ go run $(GO_FILES)

.PHONY: build
build:
	@ go build -o $(BIN_NAME) $(GO_FILES)

.PHONY: test
test:
	@ go test $(GO_PKGS)

.PHONY: clean
clean:
	@ rm $(BIN_NAME)
