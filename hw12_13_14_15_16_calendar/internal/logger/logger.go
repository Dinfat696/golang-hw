package logger

import (
    "log"
    "os"
)

type Logger struct {
    infoLog  *log.Logger
    errorLog *log.Logger
}

func New(level, location string) *Logger {
    return &Logger{
        infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
        errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
    }
}

func (l *Logger) Info(msg string) {
    l.infoLog.Println(msg)
}

func (l *Logger) Error(msg string) {
    l.errorLog.Println(msg)
}

// NewLogger - алиас для New (совместимость со старым кодом)
func NewLogger(level string) (*Logger, error) {
    return New(level, ""), nil
}
