package jsonlog

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo  Level = iota // 0.
	LevelError              // 1.
	LevelFatal              // 2.
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (logger *Logger) PrintInfo(message string, properties map[string]string) {
	_, err := logger.print(LevelInfo, message, properties)
	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}
}

func (logger *Logger) PrintError(err error, properties map[string]string) {
	_, err = logger.print(LevelError, err.Error(), properties)
	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}
}

func (logger *Logger) PrintFatal(err error, properties map[string]string) {
	_, err = logger.print(LevelFatal, err.Error(), properties)
	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}
	os.Exit(1)
}

func (logger *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < logger.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	logger.mu.Lock()
	defer logger.mu.Unlock()

	return logger.out.Write(append(line, '\n'))
}

func (logger *Logger) Write(message []byte) (n int, err error) {
	return logger.print(LevelError, string(message), nil)
}
