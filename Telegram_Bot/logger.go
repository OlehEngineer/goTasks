package main

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func StartLogging() (*os.File, error) {
	//open/create Log file, set logger setting.
	LogFile, errLog := os.OpenFile("LOG.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	//LOGRUS setting
	multiOutput := io.MultiWriter(os.Stdout, LogFile) //set logging into standard output and into a file
	log.SetOutput(multiOutput)
	log.SetLevel(NewLogger(CurrentConfiguration.LogrusLevel)) //Set log level
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	if errLog != nil {
		log.Fatalf("problem with log file, error => %s", errLog)
		return nil, errLog
	}
	return LogFile, nil
}
func NewLogger(lvl string) log.Level {
	//return logrus.SetLevel
	switch lvl {
	case "panic":
		return log.PanicLevel
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn", "warning":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	default:
		return log.InfoLevel
	}
}
