package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	LevelDebug   = "DEBUG"
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
	LevelFatal   = "FATAL"
)

var (
	mu         sync.Mutex
	warningCnt int
	errCnt     int
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
	mu.Lock()
	warningCnt++
	mu.Unlock()
	Print(LevelWarning, msg, args...)
}

func Error(msg string, args ...interface{}) {
	mu.Lock()
	errCnt++
	mu.Unlock()
	Print(LevelError, msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	Print(LevelFatal, msg, args...)
	os.Exit(1)
}

func GetWarningCount() int {
	return warningCnt
}

func GetErrorCount() int {
	return errCnt
}
