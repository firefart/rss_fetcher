@echo off
set GOOS=linux
set GOARCH=amd64

echo Updating Dependencies
go get -u gopkg.in/gomail.v2
go get -u github.com/mmcdole/gofeed

echo Building program
go build -o rss_fetcher