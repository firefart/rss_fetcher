package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
	gomail "gopkg.in/gomail.v2"
)

func sendEmail(config configuration, m *gomail.Message) error {
	debugOutput("sending mail")
	if *test {
		text, err := messageToString(m)
		if err != nil {
			return fmt.Errorf("could not print mail: %v", err)
		}
		log.Printf("[MAIL] %s", text)
		return nil
	}
	d := gomail.Dialer{Host: config.Mailserver, Port: config.Mailport}
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // nolint: gosec
	return d.DialAndSend(m)
}

func sendErrorMessage(config configuration, errorMessage error) error {
	debugOutput("sending error mail")
	m := gomail.NewMessage()
	m.SetHeader("From", config.Mailfrom)
	m.SetHeader("To", config.Mailtoerror)
	m.SetHeader("Subject", "ERROR in rss_fetcher")
	m.SetBody("text/plain", fmt.Sprintf("%v", errorMessage))

	return sendEmail(config, m)
}

func sendFeedItem(config configuration, title string, item *gofeed.Item) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Mailfrom)
	m.SetHeader("To", config.Mailto)
	m.SetHeader("Subject", fmt.Sprintf("[RSS] [%s]: %s", title, item.Title))
	m.SetBody("text/plain", feedToText(item, false))
	m.AddAlternative("text/html", feedToText(item, true))

	return sendEmail(config, m)
}

func messageToString(m *gomail.Message) (string, error) {
	buf := new(bytes.Buffer)
	_, err := m.WriteTo(buf)
	if err != nil {
		return "", fmt.Errorf("could not convert message to string: %v", err)
	}
	return buf.String(), nil
}
