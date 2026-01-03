package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stderr)

func Error(msg interface{}, keyvals ...interface{}) {
	logger.Error(msg, keyvals...)
}

func ErrorF(msg string, keyvals ...interface{}) {
	logger.Errorf(msg, keyvals...)
}

func Print(msg interface{}, keyvals ...interface{}) {
	logger.Print(msg, keyvals...)
}

func PrintF(msg string, keyvals ...interface{}) {
	logger.Printf(msg, keyvals...)
}

func Info(msg interface{}, keyvals ...interface{}) {
	logger.Info(msg, keyvals...)
}

func Warn(msg interface{}, keyvals ...interface{}) {
	logger.Warn(msg, keyvals...)
}

func Debug(msg interface{}, keyvals ...interface{}) {
	logger.Debug(msg, keyvals...)
}

func Fatal(msg interface{}, keyvals ...interface{}) {
	logger.Fatal(msg, keyvals...)
}
