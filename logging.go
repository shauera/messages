package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	config "github.com/spf13/viper"
)

// utcJSONFormatter is a logrus JSON Formatter wrapper that forces the event time to be UTC on format.
type utcJSONFormatter struct {
	fmt *log.JSONFormatter
}

// Format the log entry by forcing the event time to be UTC and delegating to the wrapped Formatter.
func (u utcJSONFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.fmt.Format(e)
}

func InitLogging() {
	// Set minimal logging level
	logLevel, err := log.ParseLevel(config.GetString("logging.level"))
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	// Log file, function and line number
	log.SetReportCaller(true)

	// Log as JSON instead of the default ASCII formatter enforcing UTC timezone
	formatter := utcJSONFormatter{fmt: new(log.JSONFormatter)}

	// Format log to show short file and function names
	gitPath := "github.com/shauera/messages"
	repoPath := fmt.Sprintf("%s/src/"+gitPath, os.Getenv("GOPATH"))
	formatter.fmt.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
		fileName := strings.Replace(f.File, repoPath, "", -1)
		functionName := strings.Replace(f.Function, gitPath, "", -1)
		return fmt.Sprintf("%s()", functionName), fmt.Sprintf("%s:%d", fileName, f.Line)
	}
	log.SetFormatter(formatter)
}
