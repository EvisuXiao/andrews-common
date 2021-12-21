package logging

import (
	"fmt"
	"log"
	"os"
)

const (
	LevelDebug   = "DEBUG"
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
	LevelFatal   = "FATAL"
)

func Print(level, msg string, args ...interface{}) {
	msg += "\n"
	msg = fmt.Sprintf("[%s] %s", level, msg)
	log.Printf(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	Print(LevelDebug, msg, args...)
}

func Info(msg string, args ...interface{}) {
	Print(LevelInfo, msg, args...)
}

func Warning(msg string, args ...interface{}) {
	Print(LevelWarning, msg, args...)
}

func Error(msg string, args ...interface{}) {
	Print(LevelError, msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	Print(LevelFatal, msg, args...)
	os.Exit(1)
}
