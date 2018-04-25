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

  "gopkg.in/gomail.v2"
  "github.com/mmcdole/gofeed"
)

const (
  timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

var (
  debug  = flag.Bool("debug", false, "Print debug output")
  config = flag.String("config", "", "Config File to use")

  configuration Configuration

  lastUpdated = time.Now()
)

type Configuration struct {
  Mailserver      string
  Mailport        int
  Mailfrom        string
  Mailtoerror     string
  Mailto          string
  Feeds           []ConfigurationFeed
  Lastupdatefile  string
}

type ConfigurationFeed struct {
  Title   string
  Url     string
}

// returns Time.now() on error
func getLastUpdated() time.Time {
  content, err := ioutil.ReadFile(configuration.Lastupdatefile)
  if err != nil {
    debugOutput(err.Error())
    return time.Now()
  }
  s, err := time.Parse(timeFormat, string(content))
  if err != nil {
    debugOutput(err.Error())
    return time.Now()
  }
  return s
}

func setLastUpdated(t time.Time) error {
  err := ioutil.WriteFile(configuration.Lastupdatefile, []byte(t.Format(timeFormat)), 0644)
  return err
}

func debugOutput(s string) {
  if *debug {
    log.Printf("[DEBUG] %s", s)
  }
}

func sendEmail(m *gomail.Message) (err error) {
  debugOutput("Sending Mail")
  d := gomail.Dialer{Host: configuration.Mailserver, Port: configuration.Mailport}
  d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
  err = d.DialAndSend(m)
  return
}

func sendErrorMessage(errorString string) error {
  m := gomail.NewMessage()
  m.SetHeader("From", configuration.Mailfrom)
  m.SetHeader("To", configuration.Mailtoerror)
  m.SetHeader("Subject", "ERROR in rss_fetcher")
  m.SetBody("text/plain", errorString)

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
  m.SetHeader("From", configuration.Mailfrom)
  m.SetHeader("To", configuration.Mailto)
  m.SetHeader("Subject", fmt.Sprintf("[RSS] [%s]: %s", title, item.Title))
  m.SetBody("text/plain", feedToText(item, false))
  m.AddAlternative("text/html", feedToText(item, true))

  err := sendEmail(m)
  return err
}

func processFeed(feedInput ConfigurationFeed) error {
  fp := gofeed.NewParser()
  feed, err := fp.ParseURL(feedInput.Url)
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

func main() {
  flag.Parse()

  if *config == "" {
    log.Fatalln("Please provide a valid config file")
  }

  file, err := os.Open(*config)
  if err != nil {
    log.Fatalf("Error opening config file: %s", err.Error())
  }
  defer file.Close()
  decoder := json.NewDecoder(file)
  configuration = Configuration{}
  err = decoder.Decode(&configuration)
  if err != nil {
    log.Fatalf("Error parsing config file: %s", err.Error())
  }

  log.Println("Starting RSS Fetcher")
  start := time.Now()
  lastUpdated = getLastUpdated()

  for _, feed := range configuration.Feeds {
    log.Printf("Processing feed %s (%s)", feed.Title, feed.Url)
    err = processFeed(feed)
    if err != nil {
      log.Printf("ERROR: %s\n", err.Error())
      err = sendErrorMessage(err.Error())
      if err != nil {
        log.Printf("ERROR on sending error mail: %s", err.Error())
      }
    }
  }
  err = setLastUpdated(start)
  if err != nil {
    log.Printf(err.Error())
  }
}
