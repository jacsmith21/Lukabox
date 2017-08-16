package logrus

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.Info("logrus init")

	log.SetLevel(log.DebugLevel)

	_, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		//log.SetOutput(file)
	} else {
		//log.Info("failed to log to file, using default")
	}

	log.Info("logrus init ending")
}

//Info Info wrapper method
func Info(args ...interface{}) {
	log.Info(args...)
}

//Debug Debug wrapper method
func Debug(args ...interface{}) {
	log.Debug(args...)
}

//Fatal Fatal wrapper method
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

//WithField WithField wrapper method
func WithField(key string, args interface{}) *log.Entry {
	return log.WithField(key, args)
}

//WithFields WithFields wrapper method
func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}
