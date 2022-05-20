package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Timeout           int                 `json:"timeout"`
	Mailserver        string              `json:"mailserver"`
	Mailport          int                 `json:"mailport"`
	Mailfrom          string              `json:"mailfrom"`
	Mailonerror       bool                `json:"mailonerror"`
	Mailtoerror       string              `json:"mailtoerror"`
	Mailto            string              `json:"mailto"`
	Feeds             []ConfigurationFeed `json:"feeds"`
	Database          string              `json:"database"`
	GlobalIgnoreWords []string            `json:"globalignorewords"`
	Test              bool
}

type ConfigurationFeed struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	IgnoreWords []string `json:"ignorewords"`
}

func GetConfig(f string) (*Configuration, error) {
	if f == "" {
		return nil, fmt.Errorf("please provide a valid config file")
	}

	b, err := os.ReadFile(f) // nolint: gosec
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b)

	decoder := json.NewDecoder(reader)
	c := Configuration{}
	if err = decoder.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
