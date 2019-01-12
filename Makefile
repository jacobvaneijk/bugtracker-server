GO_PKGS = $(shell go list ./... | sed 's/github.com\/jacobvaneijk\/bugtracker-server/./')
GO_FILES = $(shell find . -name "*.go" | grep -v "_test")

run:
	@ go run $(GO_FILES)

test:
	@ go test $(GO_PKGS)
