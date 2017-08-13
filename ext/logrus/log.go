package log

import logger "github.com/Sirupsen/logrus"

func init() {
	logger.SetLevel(logger.DebugLevel)
}

//Info Info wrapper method
func Info(args ...interface{}) {
	logger.Info(args...)
}

//Debug Debug wrapper method
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

//Fatal Fatal wrapper method
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

//WithField WithField wrapper method
func WithField(key string, args interface{}) *logger.Entry {
	return logger.WithField(key, args)
}

//WithFields WithFields wrapper method
func WithFields(fields logger.Fields) *logger.Entry {
	return logger.WithFields(fields)
}
