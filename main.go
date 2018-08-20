package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	debug      = flag.Bool("debug", false, "Print debug output")
	nomail     = flag.Bool("nomail", false, "Print mail to stdout instead of sending")
	configFile = flag.String("config", "", "Config File to use")
)

func main() {
	flag.Parse()
	config, err := getConfig(*configFile)
	if err != nil {
		log.Fatalf("could not parse config file: %v", err)
	}

	log.Println("Starting RSS Fetcher")
	start := time.Now().UnixNano()
	r, err := readDatabase(config.Database)
	if err != nil {
		processError(*config, fmt.Errorf("error in database file: %v", err))
		os.Exit(1)
	}

	cleanupDatabase(r, *config)

	for _, feed := range config.Feeds {
		log.Printf("processing feed %q (%s)", feed.Title, feed.URL)
		last, ok := r.Feeds[feed.URL]
		// if it's a new feed only process new entries and ignore old ones
		if !ok {
			last = start
		}
		entry, errFeed := processFeed(*config, feed, last)
		if errFeed != nil {
			processError(*config, errFeed)
		} else {
			r.Feeds[feed.URL] = entry
		}
	}
	r.LastRun = start
	err = saveDatabase(config.Database, r)
	if err != nil {
		processError(*config, fmt.Errorf("error on writing database file: %v", err))
		os.Exit(1)
	}
}
