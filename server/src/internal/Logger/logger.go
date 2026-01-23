package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ========= static private fields =========  //

var mu sync.Mutex = sync.Mutex{}
var logger Logger = Logger{
	logFilePath: "logs.txt",
}

// ========= structs =========  //

type Logger struct {
	logFilePath string
	file        *os.File
}

type LogEvent struct {
	Message   string
	timestamp time.Time
	severity  uint8
}

// ========= logic =========  //

func (l LogEvent) GetFormattedMessage() string {
	return fmt.Sprintf("[%s] %s-%s", l.timestamp.Format(time.RFC3339), SeverityName(l.severity), l.Message)
}

func initLogger() {
	if logger.file != nil {
		return
	}

	file, err := openFile()

	if err != nil {
		panic("Unable to initialise logger")
	}

	logger.file = file
}

func Log(message string, severity uint8) error {
	mu.Lock()
	defer mu.Unlock()

	initLogger()
	var err error

	event := LogEvent{
		Message:   message,
		timestamp: time.Now(),
		severity:  severity,
	}

	formatMessage := event.GetFormattedMessage()
	fmt.Println(formatMessage)
	_, writeErr := logger.file.WriteString(formatMessage + "\n")

	if writeErr != nil {
		panic(writeErr)
	}

	return err
}

func openFile() (*os.File, error) {
	file, err := os.OpenFile(logger.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err) // todo maybe not panic
	}
	return file, err
}

func SeverityName(s uint8) string {
	switch s {
	case 0:
		return "Info"
	case 1:
		return "Warn"
	case 2:
		return "Error"
	case 3:
		return "Fatal Error"
	default:
		return "Unknown Severity"
	}
}
