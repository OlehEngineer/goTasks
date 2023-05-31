package usecases

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func StartLogging(loggingLVL string) {
	//Logrus setting

	log.SetOutput(os.Stdout)
	log.SetLevel(LoggerLevel(loggingLVL))
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
func LoggerLevel(lvl string) log.Level {

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
