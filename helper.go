package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func processError(config configuration, err error) {
	if err != nil {
		log.Error(err)
		if config.Mailonerror {
			err = sendErrorMessage(config, err)
			if err != nil {
				log.Errorf("ERROR on sending error mail: %v", err)
			}
		}
	}
}

func timeToString(t time.Time) string {
	return t.Local().Format(time.ANSIC)
}
