package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//DefaultLogger log
var DefaultLogger *log.Entry

func init() {

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	DefaultLogger = log.WithFields(log.Fields{
		"app": "gotribe/compiler",
	})

	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
}
