package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/FireFart/rss_fetcher/internal/config"
	"github.com/FireFart/rss_fetcher/internal/database"
	"github.com/FireFart/rss_fetcher/internal/helper"
)

var (
	configFile = flag.String("config", "", "Config File to use")
)

func main() {
	flag.Parse()

	config, err := config.GetConfig(*configFile)
	if err != nil {
		log.Fatalf("could not parse config file: %v", err)
	}

	r, err := database.ReadDatabase(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Run: %s\n", helper.TimeToString(time.UnixMicro(r.LastRun)))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for key, element := range r.Feeds {
		fmt.Fprintf(w, "%s\t%s\n", key, helper.TimeToString(time.UnixMicro(element)))
	}
	w.Flush()
}
