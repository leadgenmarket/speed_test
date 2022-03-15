package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(msg string, params map[string]interface{})
	Info(msg string, params map[string]interface{})
	Warn(msg string, params map[string]interface{})
	Error(msg string, err error, params map[string]interface{})
	Panic(msg string, err error, params map[string]interface{})
}

type logger struct {
	log *logrus.Logger
}

func NewLogger(appName string, logLevel uint32, file *os.File) Logger {
	log := logrus.New()
	log.SetLevel(logrus.Level(logLevel))
	log.SetOutput(file)

	return &logger{
		log: log,
	}
}

func (lg *logger) Debug(msg string, params map[string]interface{}) {
	logger := lg.log
	for key, param := range params {
		logger.WithField(key, param)
	}
	logger.Debug(msg)
}

func (lg *logger) Info(msg string, params map[string]interface{}) {
	logger := lg.log
	for key, param := range params {
		logger.WithField(key, param)
	}
	logger.Info(msg)
}

func (lg *logger) Warn(msg string, params map[string]interface{}) {
	logger := lg.log
	for key, param := range params {
		logger.WithField(key, param)
	}
	logger.Warn(msg)
}

func (lg *logger) Error(msg string, err error, params map[string]interface{}) {
	logger := lg.log
	for key, param := range params {
		logger.WithField(key, param)
	}
	logger.WithError(errors.WithStack(err)).Error(err)
}

func (lg *logger) Panic(msg string, err error, params map[string]interface{}) {
	logger := lg.log
	for key, param := range params {
		logger.WithField(key, param)
	}
	logger.WithError(errors.WithStack(err)).Panic(err)
}
