package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Level string

const (
	LV_DEBUG = "Debug"
	LV_INFO  = "Info"
	LV_ERROR = "Error"
	LV_FATAL = "Fatal"
)

type messageInfo struct {
	Time    string `json:"time"`
	Message string `json:"message"`
	Level   Level  `json:"level"`
	Trace   string `json:"trace"`
}

func (mi messageInfo) String() string {
	return fmt.Sprintf("{\n\tMessage: %s\n\tTime: %s\n\tLevel: %s\n\tTrace: %s\n}", mi.Message, mi.Time, string(mi.Level), mi.Trace)
}

type FnCreateMessageInfo func(level Level, message, filename string, lineNumber int) messageInfo

type loggerInfo struct {
	MsgChannel    chan messageInfo
	Verbosity     Verbosity
	File          *os.File
	WaitWriteFile *sync.WaitGroup

	CreateMessage FnCreateMessageInfo
}

var logger *loggerInfo = nil

func CreateLogger(config LoggerConfiguration) (err error) {
	if config.Verbosity == MUTE {
		// We make nothing if the verbosity is muted
		return err
	}

	splitPath := strings.Split(config.LogPath, "/")
	folder := strings.Join(splitPath[:len(splitPath)-1], "/")

	if err = os.MkdirAll(folder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to open/create logger folder, with error: %w", err)
	}

	var file *os.File = nil
	if file, err = os.Create(config.LogPath); err != nil {
		return fmt.Errorf("failed to open logger file, with error: %w", err)
	}
	file.Write([]byte("[\n"))

	channel := make(chan messageInfo, int(config.MessageCapacity))
	var waitWriteFile sync.WaitGroup
	waitWriteFile.Add(1)

	go func(messageChannel chan messageInfo, wg *sync.WaitGroup) {
		for message := range messageChannel {
			if byteMessage, err := json.Marshal(message); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to marshal message %q, with error: %v", message.String(), err)

			} else if _, err = file.Write(append(byteMessage, byte('\n'))); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to log file at %q, with the message %q, with error: %v", config.LogPath, string(byteMessage), err)
			}
		}
		wg.Done()
	}(channel, &waitWriteFile)

	logger = &loggerInfo{
		MsgChannel:    channel,
		Verbosity:     config.Verbosity,
		File:          file,
		WaitWriteFile: &waitWriteFile,

		CreateMessage: func(level Level, message, filename string, lineNumber int) messageInfo {
			return messageInfo{
				Message: message,
				Time:    time.Now().Format(config.DateFormat),
				Level:   level,
				Trace:   fmt.Sprintf("In %q, at %d", filename, lineNumber),
			}
		},
	}

	return err
}

func Info(format string, args ...any) {
	if logger == nil {
		return
	}

	if _, filePath, lineNumber, ok := runtime.Caller(1); ok {
		logger.MsgChannel <- logger.CreateMessage(
			LV_INFO, fmt.Sprintf(format, args...), filePath, lineNumber,
		)
	}
}

func Debug(format string, args ...any) {
	if logger == nil || logger.Verbosity == NORMAL {
		return
	}

	if _, filePath, lineNumber, ok := runtime.Caller(1); ok {
		logger.MsgChannel <- logger.CreateMessage(
			LV_DEBUG, fmt.Sprintf(format, args...), filePath, lineNumber,
		)
	}
}

func Error(format string, args ...any) {
	if logger == nil {
		return
	}

	if _, filePath, lineNumber, ok := runtime.Caller(1); ok {
		logger.MsgChannel <- logger.CreateMessage(
			LV_ERROR, fmt.Sprintf(format, args...), filePath, lineNumber,
		)
	}
}

func Fatal(format string, args ...any) {
	if logger == nil {
		return
	}

	if _, filePath, lineNumber, ok := runtime.Caller(1); ok {
		logger.MsgChannel <- logger.CreateMessage(
			LV_FATAL, fmt.Sprintf(format, args...), filePath, lineNumber,
		)
	}
}

func Close() {
	if logger == nil {
		return
	}

	close(logger.MsgChannel)
	logger.WaitWriteFile.Wait()

	logger.File.Write([]byte("[\n"))
}
