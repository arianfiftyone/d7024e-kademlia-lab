package logger

import (
	"errors"
)

type Logger struct {
	newLogs []string
	oldLogs []string
}

func newLogger() *Logger {
	return &Logger{
		newLogs: []string{},
		oldLogs: []string{},
	}
}

var logger *Logger

func Log(log string) {
	if logger == nil {
		logger = newLogger()
	}
	logger.newLogs = append(logger.newLogs, log)

}

func ReadNewLog() (string, error) {
	if logger == nil {
		logger = newLogger()
	}

	if len(logger.newLogs) > 0 {
		oldest := logger.newLogs[0]

		logger.newLogs[0] = logger.newLogs[len(logger.newLogs)-1]
		logger.newLogs[len(logger.newLogs)-1] = ""
		logger.newLogs = logger.newLogs[:len(logger.newLogs)-1]
		logger.oldLogs = append(logger.oldLogs, oldest)

		return oldest, nil
	}
	return "", errors.New("no new logs available")
}

func GetOldLogs() []string {
	if logger == nil {
		logger = newLogger()
	}

	return logger.oldLogs
}
