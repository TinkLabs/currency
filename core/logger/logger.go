package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	appEnv := os.Getenv("APP_ENV")

	switch appEnv {
	case "testing":
		fallthrough
	case "localhost":
		// log as text form easier for human to read
		logrus.SetFormatter(&logrus.TextFormatter{})

		// log to stdout
		logrus.SetOutput(os.Stdout)

		// lower log level for testing and debugging
		logrus.SetLevel(logrus.DebugLevel)
	case "development":
		setupLogger()
	case "staging":
		setupLogger()
	case "production":
		setupLogger()

		// only log the info severity or above.
		logrus.SetLevel(logrus.InfoLevel)
	default:
		fmt.Printf("Failed to start, since APP_ENV = %s is not registered\n", appEnv)
		os.Exit(1)
	}
}

func setupLogger() {
	// log to file with json format and higher logging level
	file, err := os.OpenFile("./logs/currency.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		fmt.Printf("Failed to open log file, err=%+v\n", err)
		os.Exit(1)
	}

	// log as JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// log to file
	logrus.SetOutput(file)

	// only log the info severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}