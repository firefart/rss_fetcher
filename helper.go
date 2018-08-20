package main

import (
	"fmt"
	"log"
	"time"
)

func debugOutput(s string, a ...interface{}) {
	if *debug {
		log.Printf("[DEBUG] %s", fmt.Sprintf(s, a...))
	}
}

func processError(config configuration, err error) {
	if err != nil {
		log.Printf("ERROR: %v", err)
		if config.Mailonerror {
			err = sendErrorMessage(config, err)
			if err != nil {
				log.Printf("ERROR on sending error mail: %v", err)
			}
		}
	}
}

func timeToString(t time.Time) string {
	return t.Local().Format(time.ANSIC)
}
