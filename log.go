package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// CreateLogger can create logger ctx with fields
func CreateLogger(fields log.Fields) log.Interface {
	log.SetHandler(cli.Default)
	return log.WithFields(fields)
}

// LogSuccess can log message use info level
func LogSuccess(fields log.Fields) {
	logger := CreateLogger(fields)
	logger.Info("success")
}

// LogError can log error with error message
func LogError(errMessage string) {
	CreateLogger(log.Fields{
		"message": errMessage,
	}).Error("error")
}
