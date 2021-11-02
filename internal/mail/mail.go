package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/FireFart/rss_fetcher/internal/config"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

func SendEmail(c config.Configuration, m *gomail.Message) error {
	log.Debug("sending mail")
	if c.Test {
		text, err := messageToString(m)
		if err != nil {
			return fmt.Errorf("could not print mail: %v", err)
		}
		log.Infof("[MAIL] %s", text)
		return nil
	}
	d := gomail.Dialer{Host: c.Mailserver, Port: c.Mailport}
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // nolint: gosec
	return d.DialAndSend(m)
}

func SendErrorMessage(c config.Configuration, errorMessage error) error {
	log.Debug("sending error mail")
	m := gomail.NewMessage()
	m.SetHeader("From", c.Mailfrom)
	m.SetHeader("To", c.Mailtoerror)
	m.SetHeader("Subject", "ERROR in rss_fetcher")
	m.SetBody("text/plain", fmt.Sprintf("%v", errorMessage))

	return SendEmail(c, m)
}

func SendFeedItem(c config.Configuration, title string, item *gofeed.Item) error {
	m := gomail.NewMessage()
	m.SetHeader("From", c.Mailfrom)
	m.SetHeader("To", c.Mailto)
	m.SetHeader("Subject", fmt.Sprintf("[RSS] [%s]: %s", title, item.Title))
	m.SetBody("text/plain", feedToText(item, false))
	m.AddAlternative("text/html", feedToText(item, true))

	return SendEmail(c, m)
}

func messageToString(m *gomail.Message) (string, error) {
	buf := new(bytes.Buffer)
	_, err := m.WriteTo(buf)
	if err != nil {
		return "", fmt.Errorf("could not convert message to string: %v", err)
	}
	return buf.String(), nil
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
