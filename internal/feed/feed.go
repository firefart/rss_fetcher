package feed

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/FireFart/rss_fetcher/internal/config"
	"github.com/FireFart/rss_fetcher/internal/helper"
	"github.com/FireFart/rss_fetcher/internal/mail"

	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

func FetchFeed(url string, timeout int) (*gofeed.Feed, error) {
	t := time.Duration(timeout) * time.Second
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: t,
		}).Dial,
		TLSHandshakeTimeout: t,
	}

	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Timeout:   t,
		Transport: netTransport,
	}

	return fp.ParseURL(url)
}

func ProcessFeed(c config.Configuration, feedInput config.ConfigurationFeed, lastUpdate int64) (int64, error) {
	retVal := lastUpdate
	feed, err := FetchFeed(feedInput.URL, c.Timeout)
	if err != nil {
		return 0, fmt.Errorf("could not fetch feed %q: %v", feedInput.URL, err)
	}

	for _, item := range feed.Items {
		log.Debug(item.Title)

		if item.UpdatedParsed == nil && item.PublishedParsed == nil {
			log.Warnf("error in item for feed %s - no published or updated date", feedInput.Title)
			continue
		}

		var entryLastUpdated int64
		if item.UpdatedParsed != nil {
			entryLastUpdated = item.UpdatedParsed.UnixNano()
		} else {
			entryLastUpdated = item.PublishedParsed.UnixNano()
		}

		if entryLastUpdated > lastUpdate {
			retVal = entryLastUpdated
			log.Infof("found entry in feed %q: %q - updated: %s, lastupdated: %s", feedInput.Title, item.Title, helper.TimeToString(time.Unix(0, entryLastUpdated)), helper.TimeToString(time.Unix(0, lastUpdate)))

			words := append(c.GlobalIgnoreWords, feedInput.IgnoreWords...)
			if shouldFeedBeIgnored(words, item) {
				log.Infof("ignoring entry %q in feed %q because of matched ignore word", item.Title, feedInput.Title)
				continue
			}

			err = mail.SendFeedItem(c, feedInput.Title, item)
			if err != nil {
				return 0, err
			}
		} else {
			log.Debugf("feed %q: skipping item %q because date is in the past - updated: %s, lastupdated: %s",
				feedInput.Title, item.Title, helper.TimeToString(time.Unix(0, entryLastUpdated)), helper.TimeToString(time.Unix(0, lastUpdate)))
		}
	}

	return retVal, nil
}

func shouldFeedBeIgnored(ignoreWords []string, feed *gofeed.Item) bool {
	if helper.StringMatches(feed.Title, ignoreWords) ||
		helper.StringMatches(feed.Content, ignoreWords) ||
		helper.StringMatches(feed.Description, ignoreWords) {
		return true
	}

	return false
}
