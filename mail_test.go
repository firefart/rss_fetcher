package main

import (
	"errors"
	"testing"

	"github.com/mmcdole/gofeed"

	"gopkg.in/gomail.v2"
)

func init() {
	t := true
	test = &t
}

func TestSendEmail(t *testing.T) {
	config := configuration{}
	m := gomail.NewMessage()
	err := sendEmail(config, m)
	if err != nil {
		t.Fatalf("error returned: %v", err)
	}
}

func TestSendErrorMessage(t *testing.T) {
	config := configuration{
		Mailfrom:    "from@mail.com",
		Mailonerror: true,
		Mailtoerror: "to@mail.com",
	}
	e := errors.New("test")
	err := sendErrorMessage(config, e)
	if err != nil {
		t.Fatalf("error returned: %v", err)
	}
}
func TestSendFeedItem(t *testing.T) {
	config := configuration{
		Mailfrom:    "from@mail.com",
		Mailonerror: true,
		Mailtoerror: "to@mail.com",
	}
	i := gofeed.Item{}
	err := sendFeedItem(config, "Title", &i)
	if err != nil {
		t.Fatalf("error returned: %v", err)
	}
}
