package logger

import (
	"log"
	"os"
)

// Level controls the verbosity of logging. Supported values: "debug", "info".
var Level = os.Getenv("LOG_LEVEL")

func Debugf(format string, v ...interface{}) {
	if Level == "debug" {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if Level == "debug" || Level == "info" || Level == "" {
		log.Printf("[INFO] "+format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
