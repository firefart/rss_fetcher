@echo off

echo Updating Dependencies
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u gopkg.in/gomail.v2
go get -u github.com/mmcdole/gofeed
go get -u github.com/golang/protobuf/proto

echo Generating protobuf code
protoc --go_out=. rss/rss.proto

echo Running gometalinter
gometalinter ./...

echo Running tests
go test -v

echo Building program
set GOOS=windows
set GOARCH=amd64
go build -o rss_fetcher.exe