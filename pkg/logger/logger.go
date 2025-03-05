package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles for different log levels
	debugStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A7A7A7")).
		Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Padding(0, 1)

	warnStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Padding(0, 1)

	// File styles
	fileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Padding(0, 1)

	lineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF00FF")).
		Padding(0, 1)

	// Time style
	timeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Padding(0, 1)
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level     LogLevel
	logFile   *os.File
	debugMode bool
}

func NewLogger(debugMode bool) (*Logger, error) {
	// Create logs directory if it doesn't exist
	logDir := filepath.Join(os.Getenv("HOME"), ".tiktok-live-logger", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	logPath := filepath.Join(logDir, fmt.Sprintf("tiktok-live-%s.log", timestamp))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	level := INFO
	if debugMode {
		level = DEBUG
	}

	return &Logger{
		level:     level,
		logFile:   logFile,
		debugMode: debugMode,
	}, nil
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)

	// Format message
	formattedMsg := fmt.Sprintf(msg, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Create log entry
	logEntry := fmt.Sprintf("[%s] [%s:%d] %s\n",
		timestamp,
		file,
		line,
		formattedMsg,
	)

	// Write to file
	if l.logFile != nil {
		l.logFile.WriteString(logEntry)
	}

	// Format for console output
	var style lipgloss.Style
	switch level {
	case DEBUG:
		style = debugStyle
	case INFO:
		style = infoStyle
	case WARN:
		style = warnStyle
	case ERROR:
		style = errorStyle
	}

	// Create console output
	consoleOutput := fmt.Sprintf("%s %s %s %s",
		timeStyle.Render(timestamp),
		fileStyle.Render(fmt.Sprintf("%s:%d", file, line)),
		style.Render(fmt.Sprintf("[%s]", level.String())),
		formattedMsg,
	)

	// Print to console
	fmt.Println(consoleOutput)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

func (l *Logger) ErrorWithStack(err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	// Get stack trace
	stack := make([]byte, 4096)
	stack = stack[:runtime.Stack(stack, false)]

	// Log error with stack trace
	l.Error("%s: %v\nStack trace:\n%s", fmt.Sprintf(msg, args...), err, stack)
}

func (level LogLevel) String() string {
	switch level {
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