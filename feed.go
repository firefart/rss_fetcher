package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func fetchFeed(url string, timeout int) (*gofeed.Feed, error) {
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

func processFeed(config configuration, feedInput configurationFeed, lastUpdate int64) (int64, error) {
	retVal := lastUpdate
	feed, err := fetchFeed(feedInput.URL, config.Timeout)
	if err != nil {
		return 0, fmt.Errorf("could not fetch feed %q: %v", feedInput.URL, err)
	}

	for _, item := range feed.Items {
		debugOutput("%s", item.Title)

		if item.UpdatedParsed == nil && item.PublishedParsed == nil {
			log.Printf("error in item for feed %s - no published or updated date", feedInput.Title)
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
			log.Printf("found entry in feed %q: %q - updated: %s, lastupdated: %s", feedInput.Title, item.Title, timeToString(time.Unix(0, entryLastUpdated)), timeToString(time.Unix(0, lastUpdate)))
			err = sendFeedItem(config, feedInput.Title, item)
			if err != nil {
				return 0, err
			}
		} else {
			debugOutput("feed %q: skipping item %q because date is in the past - updated: %s, lastupdated: %s",
				feedInput.Title, item.Title, timeToString(time.Unix(0, entryLastUpdated)), timeToString(time.Unix(0, lastUpdate)))
		}
	}

	return retVal, nil
}

func feedToText(item *gofeed.Item, html bool) string {
	linebreak := "\n\n"
	if html {
		linebreak = "\n<br><br>\n"
	}
	var buffer bytes.Buffer
	if item.Link != "" {
		_, err := buffer.WriteString(fmt.Sprintf("%s%s", item.Link, linebreak))
		if err != nil {
			return err.Error()
		}
	}
	if item.Description != "" {
		_, err := buffer.WriteString(fmt.Sprintf("%s%s", item.Description, linebreak))
		if err != nil {
			return err.Error()
		}
	}
	if item.Content != "" {
		_, err := buffer.WriteString(fmt.Sprintf("%s%s", item.Content, linebreak))
		if err != nil {
			return err.Error()
		}
	}
	return strings.TrimSuffix(buffer.String(), linebreak)
}
