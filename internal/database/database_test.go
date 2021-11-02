package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/FireFart/rss_fetcher/internal/config"
)

func TestNewDatabase(t *testing.T) {
	x := NewDatabase()
	if x.Feeds == nil {
		t.Fatal("Feed map is nil")
	}
}

func TestReadEmptyDatabase(t *testing.T) {
	r, err := ReadDatabase("")
	if err != nil {
		t.Fatalf("encountered error on reading database: %v", err)
	}
	if r == nil {
		t.Fatal("returned database is nil")
	}
}

func TestReadDatabase(t *testing.T) {
	r, err := ReadDatabase(filepath.Join("testdata", "rss.testdb"))
	if err != nil {
		t.Fatalf("encountered error on reading database: %v", err)
	}
	if r == nil {
		t.Fatal("returned database is nil")
	}
}

func TestReadInvalidDatabase(t *testing.T) {
	_, err := ReadDatabase(filepath.Join("testdata", "invalid.testdb"))
	if err == nil {
		t.Fatal("expected error but none returned")
	}
}

func TestSaveDatabase(t *testing.T) {
	d := NewDatabase()
	f, err := os.CreateTemp("", "testdb")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	err = SaveDatabase(f.Name(), d)
	if err != nil {
		t.Fatalf("encountered error on writing database: %v", err)
	}
}

func TestCleanupDatase(t *testing.T) {
	initialLen := 5
	var tt = []struct {
		testName       string
		invalidEntries int
	}{
		{"No invalid entries", 0},
		{"1 invalid entry", 1},
		{"3 invalid entries", 3},
		{"300 invalid entries", 300},
	}
	for _, x := range tt {
		t.Run(x.testName, func(t *testing.T) {
			d := NewDatabase()
			c := config.Configuration{}
			for i := 0; i < initialLen; i++ {
				url := fmt.Sprintf("test%d", i)
				d.Feeds[url] = int64(i)
				x := config.ConfigurationFeed{URL: url}
				c.Feeds = append(c.Feeds, x)
			}

			// add invalid feeds to database
			for i := 0; i < x.invalidEntries; i++ {
				url := fmt.Sprintf("invalid%d", i)
				d.Feeds[url] = int64(i)
			}

			CleanupDatabase(d, c)

			if len(d.Feeds) != initialLen {
				t.Fatalf("expected Feeds to have len %d, got %d", initialLen, len(d.Feeds))
			}

			if len(d.Feeds) != len(c.Feeds) {
				t.Fatalf("expected Feeds to have same len as config %d, got %d", initialLen, len(d.Feeds))
			}
		})
	}
}
