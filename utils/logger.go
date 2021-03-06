package utils

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func InitializeLogger(fileName string) {

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	log.SetLevel(log.TraceLevel)

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	log.SetFormatter(&log.TextFormatter{})
	//log.SetFormatter(&log.JSONFormatter{})
}
