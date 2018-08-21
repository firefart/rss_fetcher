package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/FireFart/rss_fetcher/rss"
	"github.com/golang/protobuf/proto"
)

func newDatabase() *rss.Rss {
	return &rss.Rss{Feeds: make(map[string]int64)}
}

func readDatabase(database string) (*rss.Rss, error) {
	debugOutput("Reading database")
	// create database if needed
	if _, err := os.Stat(database); os.IsNotExist(err) {
		return newDatabase(), nil
	}

	b, err := ioutil.ReadFile(database) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not read database %s: %v", database, err)
	}

	var rss rss.Rss
	if err := proto.Unmarshal(b, &rss); err != nil {
		return nil, fmt.Errorf("could not unmarshal database %s: %v", database, err)
	}
	return &rss, nil
}

func saveDatabase(database string, r proto.Message) error {
	debugOutput("Saving database: %q", proto.MarshalTextString(r))
	b, err := proto.Marshal(r)
	if err != nil {
		return fmt.Errorf("could not marshal database %s: %v", database, err)
	}
	if err := ioutil.WriteFile(database, b, 0666); err != nil {
		return fmt.Errorf("could not write database %s: %v", database, err)
	}
	return nil
}

// removes old feeds from database
func cleanupDatabase(r *rss.Rss, c configuration) {
	urls := make(map[string]struct{})
	for _, x := range c.Feeds {
		urls[x.URL] = struct{}{}
	}

	for url := range r.Feeds {
		// delete entry if not present in config file
		if _, ok := urls[url]; !ok {
			debugOutput("Removing entry %q from database", url)
			delete(r.Feeds, url)
		}
	}
}
