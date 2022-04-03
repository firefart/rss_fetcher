package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/FireFart/rss_fetcher/internal/config"
	"github.com/FireFart/rss_fetcher/internal/database"
	"github.com/FireFart/rss_fetcher/internal/feed"
	"github.com/FireFart/rss_fetcher/internal/helper"

	log "github.com/sirupsen/logrus"
)

var (
	debug      = flag.Bool("debug", false, "Print debug output")
	test       = flag.Bool("test", false, "do not send mails, print them instead")
	configFile = flag.String("config", "", "Config File to use")
)

func main() {
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	config, err := config.GetConfig(*configFile)
	if err != nil {
		log.Fatalf("could not parse config file: %v", err)
	}
	config.Test = *test

	log.Println("Starting RSS Fetcher")
	start := time.Now().UnixNano()
	r, err := database.ReadDatabase(config.Database)
	if err != nil {
		helper.ProcessError(*config, fmt.Errorf("error in database file: %v", err))
		os.Exit(1)
	}

	database.CleanupDatabase(r, *config)

	for _, f := range config.Feeds {
		log.Printf("processing feed %q (%s)", f.Title, f.URL)
		last, ok := r.Feeds[f.URL]
		// if it's a new feed only process new entries and ignore old ones
		if !ok {
			last = start
		}
		entry, errFeed := feed.ProcessFeed(*config, f, last)
		if errFeed != nil {
			helper.ProcessError(*config, errFeed)
		} else {
			r.Feeds[f.URL] = entry
		}
	}
	r.LastRun = start
	err = database.SaveDatabase(config.Database, r)
	if err != nil {
		helper.ProcessError(*config, fmt.Errorf("error on writing database file: %v", err))
		os.Exit(1)
	}
}
