.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -trimpath .

.PHONY: linux
linux: protoc update test
	GOOS=linux GOARCH=amd64 go build -trimpath .

.PHONY: test
test:
	go test -v -race ./...

.PHONY: update
update: protoc
	go get -u
	go mod tidy -v

.PHONY: lint
lint:
	"$$(go env GOPATH)/bin/golangci-lint" run ./...
	go mod tidy

.PHONY: lint-update
lint-update:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint --version

.PHONY: lint-docker
lint-docker:
	docker pull golangci/golangci-lint:latest
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run

.PHONY: protoc
protoc:
	# make sure to always have the latest protoc installed with all the includes
	# https://github.com/protocolbuffers/protobuf/releases/latest
	# sudo mv bin/protoc /usr/bin/protoc
	# sudo rm -rf /usr/local/include/google
	# sudo mv include/* /usr/local/include
	# rm -rf bin
	# rm -rf include
	# rm -f readme.txt
	# get dependencies
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	protoc -I ./proto -I /usr/local/include/ ./proto/rss.proto --go_out=.
