package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type configuration struct {
	Timeout     int                 `json:"timeout"`
	Mailserver  string              `json:"mailserver"`
	Mailport    int                 `json:"mailport"`
	Mailfrom    string              `json:"mailfrom"`
	Mailonerror bool                `json:"mailonerror"`
	Mailtoerror string              `json:"mailtoerror"`
	Mailto      string              `json:"mailto"`
	Feeds       []configurationFeed `json:"feeds"`
	Database    string              `json:"database"`
}

type configurationFeed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func getConfig(f string) (*configuration, error) {
	if f == "" {
		return nil, fmt.Errorf("please provide a valid config file")
	}

	b, err := os.ReadFile(f) // nolint: gosec
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b)

	decoder := json.NewDecoder(reader)
	c := configuration{}
	if err = decoder.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
