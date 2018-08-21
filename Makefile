GOPATH := $(or $(GOPATH), $(HOME)/go)
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

.DEFAULT_GOAL := build

.PHONY: build
build: protoc deps test
	go build .

.PHONY: test
test: deps lint
	go test -v ./...

.PHONY: deps
deps:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u gopkg.in/gomail.v2
	go get -u github.com/mmcdole/gofeed
	go get -u github.com/golang/protobuf/proto

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install &> /dev/null

.PHONY: lint
lint: deps $(GOMETALINTER)
	$(BIN_DIR)/gometalinter ./...

.PHONY: protoc
protoc: deps
	protoc --go_out=. rss/rss.proto
