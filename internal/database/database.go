package database

import (
	"fmt"
	"os"

	"github.com/FireFart/rss_fetcher/internal/config"
	"github.com/FireFart/rss_fetcher/internal/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func NewDatabase() *pb.Rss {
	return &pb.Rss{Feeds: make(map[string]int64)}
}

func ReadDatabase(database string) (*pb.Rss, error) {
	log.Debug("Reading database")
	// create database if needed
	if _, err := os.Stat(database); os.IsNotExist(err) {
		return NewDatabase(), nil
	}

	b, err := os.ReadFile(database) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not read database %s: %v", database, err)
	}

	rssMsg := &pb.Rss{}
	if err := proto.Unmarshal(b, rssMsg); err != nil {
		return nil, fmt.Errorf("could not unmarshal database %s: %v", database, err)
	}
	return rssMsg, nil
}

func SaveDatabase(database string, r proto.Message) error {
	b, err := proto.Marshal(r)
	if err != nil {
		return fmt.Errorf("could not marshal database %s: %v", database, err)
	}
	if err := os.WriteFile(database, b, 0666); err != nil {
		return fmt.Errorf("could not write database %s: %v", database, err)
	}
	return nil
}

// removes old feeds from database
func CleanupDatabase(r *pb.Rss, c config.Configuration) {
	urls := make(map[string]struct{})
	for _, x := range c.Feeds {
		urls[x.URL] = struct{}{}
	}

	for url := range r.Feeds {
		// delete entry if not present in config file
		if _, ok := urls[url]; !ok {
			log.Debugf("Removing entry %q from database", url)
			delete(r.Feeds, url)
		}
	}
}
