package global

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// FileLogger implements a concurrent-safe file logger
type FileLogger struct {
	path      string
	logger    io.WriteCloser
	mu        sync.Mutex
	minLevel  LogLevel
	timestamp string
}

// NewFileLogger creates a new FileLogger instance
func NewFileLogger(path string, minLevel LogLevel) (*FileLogger, error) {
	if path == "" {
		path = "."
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	filename := filepath.Join(path, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create or open log file: %w", err)
	}

	return &FileLogger{
		path:      path,
		logger:    file,
		minLevel:  minLevel,
		timestamp: "2006-01-02 15:04:05.000",
	}, nil
}

// Close properly closes the log file
func (l *FileLogger) Close() error {
	return l.logger.Close()
}

// SetTimeFormat allows customizing the timestamp format
func (l *FileLogger) SetTimeFormat(format string) {
	l.timestamp = format
}

// log handles the actual logging with proper synchronization
func (l *FileLogger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format(l.timestamp)
	logMsg := fmt.Sprintf(msg, args...)
	row := fmt.Sprintf("%s  %-5s  %s\n", timestamp, level, logMsg)

	if _, err := fmt.Fprint(l.logger, row); err != nil {
		// write to stderr as a fallback
		fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
		fmt.Fprint(os.Stderr, row)
	}
}

// Debug logs a debug message
func (l *FileLogger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

// Info logs an info message
func (l *FileLogger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

// Warn logs a warning message
func (l *FileLogger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

// Error logs an error message
func (l *FileLogger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

// Initialize the default logger
var Log, _ = NewFileLogger(os.Getenv("LOG_DIR"), INFO)
