@echo off
set GOOS=linux
set GOARCH=amd64

echo Updating Dependencies
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u gopkg.in/gomail.v2
go get -u github.com/mmcdole/gofeed
go get -u github.com/golang/protobuf/proto

echo Running gometalinter
gometalinter ./...

echo Building program
protoc --go_out=. rss/rss.proto
go build -o rss_fetcher