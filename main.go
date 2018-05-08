package main

import (
  "encoding/json"
  "crypto/tls"
  "flag"
  "fmt"
  "log"
  "os"
  "time"
  "bytes"
  "io/ioutil"
  "strings"
  "net"
  "net/http"

  "gopkg.in/gomail.v2"
  "github.com/mmcdole/gofeed"
)

const (
  timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

var (
  debug      = flag.Bool("debug", false, "Print debug output")
  configFile = flag.String("config", "", "Config File to use")

  config configuration

  lastUpdated = time.Now()
)

type configuration struct {
  Timeout         int                   `json:"timeout"`
  Mailserver      string                `json:"mailserver"`
  Mailport        int                   `json:"mailport"`
  Mailfrom        string                `json:"mailfrom"`
  Mailonerror     bool                  `json:"mailonerror"`
  Mailtoerror     string                `json:"mailtoerror"`
  Mailto          string                `json:"mailto"`
  Feeds           []configurationFeed   `json:"feeds"`
  Lastupdatefile  string                `json:"lastupdatefile"`
}

type configurationFeed struct {
  Title   string    `json:"title"`
  URL     string    `json:"url"`
}

// returns Time.now() on error
func getLastUpdated() time.Time {
  content, err := ioutil.ReadFile(config.Lastupdatefile)
  if err != nil {
    debugOutput(fmt.Sprintf("error on reading last udpate file: %v", err))
    return time.Now()
  }
  s, err := time.Parse(timeFormat, string(content))
  if err != nil {
    debugOutput(fmt.Sprintf("error on parsing last udpate file: %v", err))
    return time.Now()
  }
  return s
}

func setLastUpdated(t time.Time) error {
  err := ioutil.WriteFile(config.Lastupdatefile, []byte(t.Format(timeFormat)), 0644)
  return err
}

func debugOutput(s string) {
  if *debug {
    log.Printf("[DEBUG] %s", s)
  }
}

func sendEmail(m *gomail.Message) (err error) {
  debugOutput("Sending Mail")
  d := gomail.Dialer{Host: config.Mailserver, Port: config.Mailport}
  d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // nolint: gas
  err = d.DialAndSend(m)
  return
}

func sendErrorMessage(errorMessage error) error {
  m := gomail.NewMessage()
  m.SetHeader("From", config.Mailfrom)
  m.SetHeader("To", config.Mailtoerror)
  m.SetHeader("Subject", "ERROR in rss_fetcher")
  m.SetBody("text/plain", fmt.Sprintf("%v", errorMessage))

  err := sendEmail(m)
  return err
}

func feedToText(item *gofeed.Item, html bool) string {
  linebreak := "\n\n"
  if html {
    linebreak = "\n<br><br>\n"
  }
  var buffer bytes.Buffer
  if item.Link != "" {
    buffer.WriteString(fmt.Sprintf("%s%s", item.Link, linebreak))
  }
  if item.Description != "" {
    buffer.WriteString(fmt.Sprintf("%s%s", item.Description, linebreak))
  }
  if item.Content != "" {
    buffer.WriteString(fmt.Sprintf("%s%s", item.Content, linebreak))
  }
  return strings.TrimSuffix(buffer.String(), linebreak)
}

func sendFeedItem(title string, item *gofeed.Item) error {
  m := gomail.NewMessage()
  m.SetHeader("From", config.Mailfrom)
  m.SetHeader("To", config.Mailto)
  m.SetHeader("Subject", fmt.Sprintf("[RSS] [%s]: %s", title, item.Title))
  m.SetBody("text/plain", feedToText(item, false))
  m.AddAlternative("text/html", feedToText(item, true))

  err := sendEmail(m)
  return err
}

// nolint: gocyclo
func processFeed(feedInput configurationFeed) error {
  timeout := time.Duration(config.Timeout) * time.Second
  netTransport := &http.Transport {
    Dial: (&net.Dialer{
      Timeout: timeout,
    }).Dial,
    TLSHandshakeTimeout: timeout,
  }

  fp := gofeed.NewParser()
  fp.Client = &http.Client{
    Timeout: timeout,
    Transport: netTransport,
  }

  feed, err := fp.ParseURL(feedInput.URL)
  if err != nil {
    return err
  }
  // check if feed was updated
  if feed.UpdatedParsed != nil && feed.UpdatedParsed.Before(lastUpdated) {
    return nil
  }

  for _, item := range feed.Items {
    if item.UpdatedParsed == nil && item.PublishedParsed == nil {
      log.Printf("Error in item for feed %s. No published or updated date", feedInput.Title)
      continue
    }
    if (item.UpdatedParsed != nil && item.UpdatedParsed.After(lastUpdated)) ||
       (item.PublishedParsed != nil && item.PublishedParsed.After(lastUpdated)) {
      debugOutput(fmt.Sprintf("Found Entry in feed %s: %s", feedInput.Title, item.Title))
      err = sendFeedItem(feedInput.Title, item)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

// nolint: gocyclo
func main() {
  flag.Parse()

  if *configFile == "" {
    log.Fatalln("Please provide a valid config file")
  }

  file, err := os.Open(*configFile)
  if err != nil {
    log.Fatalf("Error opening config file: %v", err)
  }

  defer func() {
    rerr := file.Close()
    if rerr != nil {
      log.Fatalf("Error closing config file: %v", rerr)
    }
  }()

  decoder := json.NewDecoder(file)
  config = configuration{}
  err = decoder.Decode(&config)
  if err != nil {
    log.Fatalf("Error parsing config file: %v", err)
  }

  log.Println("Starting RSS Fetcher")
  start := time.Now()
  lastUpdated = getLastUpdated()

  for _, feed := range config.Feeds {
    log.Printf("Processing feed %s (%s)", feed.Title, feed.URL)
    err = processFeed(feed)
    if err != nil {
      log.Printf("ERROR: %v\n", err)
      if config.Mailonerror {
        err = sendErrorMessage(err)
        if err != nil {
          log.Printf("ERROR on sending error mail: %v", err)
        }
      }
    }
  }
  err = setLastUpdated(start)
  if err != nil {
    log.Printf("error on writing last udpate file: %v", err)
  }
}
