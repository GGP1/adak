package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	// Log is the logger used to display all the information
	Log = New()
)

const (
	// Info designates informational messages that highlight the progress of the application at coarse-grained level.
	Info Level = iota
	// Debug designates fine-grained informational events useful to debug an application.
	Debug
	// Error designates error events.
	Error
	// Fatal shows an error and exits.
	Fatal
)

// Level represents the logging level used.
type Level uint8

// Logger contains the logging options.
type Logger struct {
	Out           io.Writer
	Prefix        string
	Level         Level
	Time          string
	ShowTimestamp bool
}

// New creates a new logger.
func New() *Logger {
	return &Logger{
		Out:           os.Stderr,
		Prefix:        "[ADAK]",
		ShowTimestamp: true,
	}
}

func (l *Logger) log(level Level, message string) {
	var lvl string

	if l.ShowTimestamp {
		l.Time = time.Now().Format("15:04:05 02/01/2006") + " "
	} else {
		l.Time = ""
	}

	switch level {
	case Info:
		lvl = "INFO"
	case Debug:
		lvl = "DEBUG"
	case Error:
		lvl = "ERROR"
	case Fatal:
		lvl = "FATAL"
	}

	log := fmt.Sprintf("%s%s - %s: %s", l.Time, l.Prefix, lvl, message)

	fmt.Fprintln(l.Out, log)
}

// Info provides useful information about the server.
func (l *Logger) Info(args ...interface{}) {
	l.log(Info, fmt.Sprint(args...))
}

// Infof is like Info but takes a formatted message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(Info, fmt.Sprintf(format, args...))
}

// Debug provides useful information for debugging.
func (l *Logger) Debug(args ...interface{}) {
	l.log(Debug, fmt.Sprint(args...))
}

// Debugf is like Debug but takes a formatted message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(Debug, fmt.Sprintf(format, args...))
}

// Error reports the application errors.
func (l *Logger) Error(args ...interface{}) {
	l.log(Error, fmt.Sprint(args...))
}

// Errorf is like Error but takes a formatted message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(Error, fmt.Sprintf(format, args...))
}

// Fatal reports the application errors and exists.
func (l *Logger) Fatal(args ...interface{}) {
	l.log(Fatal, fmt.Sprint(args...))
	os.Exit(1)
}

// Fatalf is like Fatal but takes a formatted message.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(Fatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}
