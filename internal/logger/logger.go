package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	INFO    = "INFO"
	ERROR   = "ERROR"
	DEBUG   = "DEBUG"
	WARNING = "WARNING"
)

type Logger struct {
	file *os.File
}

func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %v", err)
	}
	return &Logger{file: file}, nil
}

func (l *Logger) Close() {
	l.file.Close()
}

func (l *Logger) log(level string, msg string, file string, line int) {
	now := time.Now().Format("2006-01-02 15:04:05")

	fileName := file[strings.LastIndex(file, "/")+1:]

	logMessage := fmt.Sprintf("[%s] [%s] %s (File: %s, Line: %d)\n", now, level, msg, fileName, line)

	_, err := l.file.WriteString(logMessage)
	if err != nil {
		fmt.Println("Failed to write log to file:", err)
	}

	fmt.Print(logMessage)
}

func (l *Logger) Info(msg string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(INFO, msg, file, line)
}

func (l *Logger) Error(msg string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(ERROR, msg, file, line)
}

func (l *Logger) Debug(msg string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(DEBUG, msg, file, line)
}

func (l *Logger) Warning(msg string) {
	_, file, line, _ := runtime.Caller(1)
	l.log(WARNING, msg, file, line)
}
