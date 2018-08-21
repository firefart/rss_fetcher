GOPATH := $(or $(GOPATH), $(HOME)/go)

.DEFAULT_GOAL := build

.PHONY: build
build: protoc deps test
	go build .

.PHONY: test
test: deps lint
	go test -v ./...

.PHONY: deps
deps:
	go get -u gopkg.in/gomail.v2
	go get -u github.com/mmcdole/gofeed
	go get -u github.com/golang/protobuf/proto

.PHONY: lint
lint: deps
	go get -u github.com/alecthomas/gometalinter
	gometalinter --deadline=5m  ./...

.PHONY: protoc
protoc:
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc --go_out=. rss/rss.proto
