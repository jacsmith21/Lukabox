package log

import (
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.Info("logrus init")

	log.SetLevel(log.DebugLevel)

	/*_, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
	  log.SetOutput(file)
	} else {
	  log.Info("failed to log to file, using default")
	}*/

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

//Debugf Debugf wrapper method
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

//Fatal Fatal wrapper method
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

//Error Error wrapper method
func Error(args ...interface{}) {
	log.Error(args...)
}

//Errorf Errorf wrapper method
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

//WithField WithField wrapper method
func WithField(key string, args interface{}) *log.Entry {
	return log.WithField(key, args)
}

//WithFields WithFields wrapper method
func WithFields(fields log.Fields) *log.Entry {
	return log.WithFields(fields)
}

//WithError WithError wrapper method
func WithError(err error) *log.Entry {
	return log.WithError(err)
}
