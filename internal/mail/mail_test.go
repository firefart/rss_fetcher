package mail

import (
	"errors"
	"strings"
	"testing"

	"github.com/FireFart/rss_fetcher/internal/config"
	"github.com/mmcdole/gofeed"

	"gopkg.in/gomail.v2"
)

func TestSendEmail(t *testing.T) {
	config := config.Configuration{Test: true}
	m := gomail.NewMessage()
	err := SendEmail(config, m)
	if err != nil {
		t.Fatalf("error returned: %v", err)
	}
}

func TestSendErrorMessage(t *testing.T) {
	config := config.Configuration{
		Test:        true,
		Mailfrom:    "from@mail.com",
		Mailonerror: true,
		Mailtoerror: "to@mail.com",
	}
	e := errors.New("test")
	err := SendErrorMessage(config, e)
	if err != nil {
		t.Fatalf("error returned: %v", err)
	}
}
func TestSendFeedItem(t *testing.T) {
	config := config.Configuration{
		Test:        true,
		Mailfrom:    "from@mail.com",
		Mailonerror: true,
		Mailtoerror: "to@mail.com",
	}
	i := gofeed.Item{}
	err := SendFeedItem(config, "Title", &i)
	if err != nil {
		t.Fatalf("error returned: %v", err)
	}
}

func TestFeedToText(t *testing.T) {
	item := gofeed.Item{}
	item.Description = "Description"
	item.Link = "Link"
	item.Content = "Content"

	x := feedToText(&item, false)
	if !strings.Contains(x, "Description") {
		t.Fatal("missing description in feed text")
	}
	if !strings.Contains(x, "Link") {
		t.Fatal("missing link in feed text")
	}
	if !strings.Contains(x, "Content") {
		t.Fatal("missing content in feed text")
	}

	x = feedToText(&item, true)
	if !strings.Contains(x, "<br><br>") {
		t.Fatal("missing html line break in feed text")
	}
}
