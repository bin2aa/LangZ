package logger

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	debug *log.Logger
}

var App *Logger

func Init() {
	App = &Logger{
		info:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile),
		err:   log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		debug: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.warn.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.err.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
	App.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	App.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	App.Error(format, v...)
}

func Debug(format string, v ...interface{}) {
	App.Debug(format, v...)
}
