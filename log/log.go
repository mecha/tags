package log

import (
	"fmt"
	"os"
)

var level int = ErrorLevel

const (
	QuietLevel int = iota
	ErrorLevel
	InfoLevel
	DebugLevel
)

func SetLevel(newLvl int) {
	level = newLvl
}

func Error(msg string, args ...interface{}) {
	if level >= ErrorLevel {
		fmt.Fprintf(os.Stderr, msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	if level >= InfoLevel {
		fmt.Fprintf(os.Stderr, msg, args...)
	}
}

func Debug(msg string, args ...interface{}) {
	if level >= DebugLevel {
		fmt.Fprintf(os.Stderr, msg, args...)
	}
}
