package feed

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/FireFart/rss_fetcher/internal/config"
)

func feedServer(t *testing.T, filename string) *httptest.Server {
	t.Helper()
	fullName := filepath.Join("testdata", filename)
	b, err := os.ReadFile(fullName)
	if err != nil {
		t.Fatalf("could not read file %s: %v", fullName, err)
	}
	content := string(b)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, content)
	}))
	return ts
}

func TestFetchFeed(t *testing.T) {
	ts := feedServer(t, "valid_feed.xml")
	defer ts.Close()
	_, err := FetchFeed(ts.URL, 10)
	if err != nil {
		t.Fatalf("got error when fetching feed: %v", err)
	}
}

func TestFetchFeedInvalid(t *testing.T) {
	ts := feedServer(t, "invalid_feed.xml")
	defer ts.Close()
	_, err := FetchFeed(ts.URL, 10)
	if err == nil {
		t.Fatalf("expected error but got none")
	}
}

func TestProcessFeed(t *testing.T) {
	ts := feedServer(t, "valid_feed.xml")
	defer ts.Close()
	c := config.Configuration{
		Test:    true,
		Timeout: 1,
	}
	input := config.ConfigurationFeed{
		Title: "Title",
		URL:   ts.URL,
	}

	// with mail
	_, err := ProcessFeed(c, input, 0)
	if err != nil {
		t.Fatalf("got error when fetching feed: %v", err)
	}

	// no mail
	_, err = ProcessFeed(c, input, math.MaxInt64)
	if err != nil {
		t.Fatalf("got error when fetching feed: %v", err)
	}
}
