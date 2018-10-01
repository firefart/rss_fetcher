@echo off
SET BUILDARGS=-ldflags="-s -w" -gcflags="all=-trimpath=%GOPATH%\src" -asmflags="all=-trimpath=%GOPATH%\src"

echo Updating Dependencies
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u gopkg.in/gomail.v2
go get -u github.com/mmcdole/gofeed
go get -u github.com/golang/protobuf/proto

echo Generating protobuf code
protoc --go_out=. rss/rss.proto

echo Running gometalinter
go get -u github.com/alecthomas/gometalinter
gometalinter --install > nul
gometalinter ./...

echo Running tests
go test -v ./...

echo Running build
set GOOS=windows
set GOARCH=amd64
go build %BUILDARGS% -o rss_fetcher.exe
