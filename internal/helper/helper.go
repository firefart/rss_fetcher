package helper

import (
	"time"

	"github.com/FireFart/rss_fetcher/internal/config"
	"github.com/FireFart/rss_fetcher/internal/mail"
	log "github.com/sirupsen/logrus"
)

func ProcessError(c config.Configuration, err error) {
	if err != nil {
		log.Error(err)
		if c.Mailonerror {
			err = mail.SendErrorMessage(c, err)
			if err != nil {
				log.Errorf("ERROR on sending error mail: %v", err)
			}
		}
	}
}

func TimeToString(t time.Time) string {
	return t.Local().Format(time.ANSIC)
}
