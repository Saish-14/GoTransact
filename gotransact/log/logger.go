// logger/logger.go
package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	InfoLogger  *logrus.Logger
	ErrorLogger *logrus.Logger
)

func Init() {
	// Create the info logger
	InfoLogger = logrus.New()
	infoFile, err := os.OpenFile("log/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		InfoLogger.Out = infoFile
	} else {
		InfoLogger.Info("Failed to log to file, using default stderr")
	}

	// Create the error logger
	ErrorLogger = logrus.New()
	errorFile, err := os.OpenFile("log/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		ErrorLogger.Out = errorFile
	} else {
		ErrorLogger.Error("Failed to log to file, using default stderr")
	}

	InfoLogger.SetFormatter(&logrus.JSONFormatter{})
	ErrorLogger.SetFormatter(&logrus.JSONFormatter{})

	InfoLogger.SetReportCaller(true)
	ErrorLogger.SetReportCaller(true)

	// Set log level
	InfoLogger.SetLevel(logrus.InfoLevel)
	ErrorLogger.SetLevel(logrus.ErrorLevel)
}
