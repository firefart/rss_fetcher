BINARY := rss_fetcher
VERSION ?= latest
OUTDIR := out
ARCH := amd64
PLATFORMS := windows linux darwin
os = $(word 1, $@)

GOPATH := $(or $(GOPATH), $(HOME)/go)
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

.DEFAULT_GOAL := build

.PHONY: $(PLATFORMS)
$(PLATFORMS): protoc
	GOOS=$(os) GOARCH=$(ARCH) go build -o $(OUTDIR)/$(BINARY)-$(VERSION)-$(os)-$(ARCH)
	zip -j $(OUTDIR)/$(BINARY)-$(VERSION)-$(os)-$(ARCH).zip $(OUTDIR)/$(BINARY)-$(VERSION)-$(os)-$(ARCH)

.PHONY: release
release: clean deps lint test $(PLATFORMS)

.PHONY: build
build: protoc test
	go build .

.PHONY: test
test: protoc
	go test -v

.PHONY: deps
deps:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u gopkg.in/gomail.v2
	go get -u github.com/mmcdole/gofeed
	go get -u github.com/golang/protobuf/proto

.PHONY: clean
clean:
	rm -rf $(OUTDIR)/*

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install &> /dev/null

.PHONY: lint
lint: deps $(GOMETALINTER)
	$(BIN_DIR)/gometalinter ./...

.PHONY: protoc
protoc:
	protoc --go_out=. rss/rss.proto
