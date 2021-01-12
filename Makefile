GOPATH := $(or $(GOPATH), $(HOME)/go)

.DEFAULT_GOAL := build

.PHONY: build
build: protoc update test
	go build -trimpath .

.PHONY: linux
linux: protoc update test
	GOOS=linux GOARCH=amd64 GO111MODULE=on go build -trimpath .

.PHONY: test
test: update
	go test -v -race ./...

.PHONY: update
update:
	go get -u
	go mod tidy -v

.PHONY: lint
lint: update
	@if [ ! -f "$$(go env GOPATH)/bin/golangci-lint" ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.34.1; \
	fi
	golangci-lint run ./...
	go mod tidy

.PHONY: protoc
protoc:
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc --go_out=. rss/rss.proto
