package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
)

func feedServer(t *testing.T, filename string) *httptest.Server {
	t.Helper()
	fullName := filepath.Join("testdata", filename)
	b, err := ioutil.ReadFile(fullName)
	if err != nil {
		t.Fatalf("could not read file %s: %v", fullName, err)
	}
	content := string(b)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, content)
	}))
	return ts
}

func TestFetchFeed(t *testing.T) {
	t.Parallel()
	ts := feedServer(t, "valid_feed.xml")
	defer ts.Close()
	_, err := fetchFeed(ts.URL, 10)
	if err != nil {
		t.Fatalf("got error when fetching feed: %v", err)
	}
}

func TestFetchFeedInvalid(t *testing.T) {
	t.Parallel()
	ts := feedServer(t, "invalid_feed.xml")
	defer ts.Close()
	_, err := fetchFeed(ts.URL, 10)
	if err == nil {
		t.Fatalf("expected error but got none")
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