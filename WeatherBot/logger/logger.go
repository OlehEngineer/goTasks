package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func StartLogging(loglvl string) {

	//LOGRUS setting
	log.SetOutput(os.Stdout)        //set logging into standard output
	log.SetLevel(NewLogger(loglvl)) //Set log level
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

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
