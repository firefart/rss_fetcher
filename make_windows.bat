@echo off

SET BUILDARGS=-ldflags="-s -w" -trimpath

echo [*] Updating Dependencies
go get -u

echo [*] mod tidy
go mod tidy -v

echo [*] Generating protobuf code
go get -u github.com/golang/protobuf/protoc-gen-go
protoc --go_out=. rss/rss.proto

echo [*] mod tidy v2
go mod tidy -v

echo [*] go fmt
go fmt ./...

echo [*] go vet
go vet ./...

echo [*] Linting
go get -u github.com/golangci/golangci-lint@master
golangci-lint run ./...
go mod tidy

echo [*] Running Tests
go test -v ./...

echo [*] Running build
set GOOS=windows
set GOARCH=amd64
go build %BUILDARGS% -o rss_fetcher.exe
