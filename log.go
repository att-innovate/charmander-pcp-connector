package pcp

import (
	"errors"
	"log"
)

const (
	LOG_DISABLE = 1 << iota
	LOG_ERRORS
	LOG_WARNINGS
	LOG_INFO
	LOG_DEBUG
)

type Logger struct {
	LogLevel int
}

func NewLogger(loglevel int) *Logger {
	level := LOG_DEBUG

	switch {
	case loglevel == LOG_ERRORS:
		level = LOG_ERRORS
	case loglevel == LOG_WARNINGS:
		level = LOG_ERRORS | LOG_WARNINGS
	case loglevel == LOG_INFO:
		level = LOG_ERRORS | LOG_WARNINGS | LOG_INFO
	case loglevel == LOG_DEBUG:
		level = LOG_ERRORS | LOG_WARNINGS | LOG_INFO | LOG_DEBUG
	}
	return &Logger{LogLevel: level}
}

func (l *Logger) SetLogLevel(level int) error {
	if level&(LOG_INFO|
		LOG_DEBUG|
		LOG_ERRORS|
		LOG_WARNINGS|
		LOG_DISABLE) == 0 {
		return errors.New("Log level not recognized!")
	}
	l.LogLevel = level
	return nil
}

func (l *Logger) Debugln(line interface{}) {
	if (l.LogLevel & LOG_DEBUG) == LOG_DEBUG {
		log.Println(line)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if (l.LogLevel & LOG_DEBUG) == LOG_DEBUG {
		log.Printf(format, args)
	}
}

func (l *Logger) Infoln(line interface{}) {
	if (l.LogLevel & LOG_INFO) == LOG_INFO {
		log.Println(line)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if (l.LogLevel & LOG_INFO) == LOG_INFO {
		log.Printf(format, args)
	}
}

func (l *Logger) Errorln(line string) {
	if (l.LogLevel & LOG_ERRORS) == LOG_ERRORS {
		log.Printf(line)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if (l.LogLevel & LOG_ERRORS) == LOG_ERRORS {
		log.Printf(format, args)
	}
}
