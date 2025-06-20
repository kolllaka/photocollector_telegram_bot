package logging

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogFileMaxSizeMB  = 10
	defaultLogFileMaxBackups = 3
	defaultLogFileMaxAgeDays = 14
)

type LoggerFileOptions struct {
	LogFilePath       string
	LogFileMaxSizeMB  int
	LogFileMaxBackups int
	LogFileMaxAgeDays int
	LogFileCompress   bool
}

type LoggerFileOption func(*LoggerFileOptions)

func NewFileWriter(filePath string, opts ...LoggerFileOption) io.Writer {
	var w io.Writer
	config := &LoggerFileOptions{
		LogFilePath:       filePath,
		LogFileMaxSizeMB:  defaultLogFileMaxSizeMB,
		LogFileMaxBackups: defaultLogFileMaxBackups,
		LogFileMaxAgeDays: defaultLogFileMaxAgeDays,
	}

	for _, opt := range opts {
		opt(config)
	}

	w = &lumberjack.Logger{
		Filename:   config.LogFilePath,
		MaxSize:    config.LogFileMaxSizeMB,
		MaxBackups: config.LogFileMaxBackups,
		MaxAge:     config.LogFileMaxAgeDays,
		Compress:   config.LogFileCompress,
	}

	return w
}

// WithLogFilePath logger option sets the file where logs will be written.
func WithLogFilePath(logFilePath string) LoggerFileOption {
	return func(o *LoggerFileOptions) {
		o.LogFilePath = logFilePath
	}
}

// WithLogFileMaxSizeMB logger option sets the maximum file size for rotation.
func WithLogFileMaxSizeMB(maxSize int) LoggerFileOption {
	return func(o *LoggerFileOptions) {
		o.LogFileMaxSizeMB = maxSize
	}
}

// WithLogFileMaxBackups logger option sets the number of backup files to retain.
func WithLogFileMaxBackups(maxBackups int) LoggerFileOption {
	return func(o *LoggerFileOptions) {
		o.LogFileMaxBackups = maxBackups
	}
}

// WithLogFileMaxAgeDays logger option sets the maximum age of the log files.
func WithLogFileMaxAgeDays(maxAge int) LoggerFileOption {
	return func(o *LoggerFileOptions) {
		o.LogFileMaxAgeDays = maxAge
	}
}

// WithLogFileCompress logger options set needs compression.
func WithLogFileCompress(compression bool) LoggerFileOption {
	return func(o *LoggerFileOptions) {
		o.LogFileCompress = compression
	}
}
