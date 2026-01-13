package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Verbosity string

const (
	MUTE    = "mute"
	NORMAL  = "normal"
	VERBOSE = "verbose"
)

type loggerInfo struct {
	MsgChannel    chan string
	Verbosity     Verbosity
	File          *os.File
	WaitWriteFile *sync.WaitGroup
}

var logger *loggerInfo = nil

func CreateLogger(path string, verbosity Verbosity, messageCapacity uint) (err error) {
	if verbosity == MUTE {
		// We make nothing if the verbosity is muted
		return err
	}

	splitPath := strings.Split(path, "/")
	folder := strings.Join(splitPath[:len(splitPath)-1], "/")

	if err = os.MkdirAll(folder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to open/create logger folder, with error: %w", err)
	}

	var file *os.File = nil
	if file, err = os.Create(path); err != nil {
		return fmt.Errorf("failed to open logger file, with error: %w", err)
	}

	channel := make(chan string, messageCapacity)
	var waitWriteFile sync.WaitGroup
	waitWriteFile.Add(1)

	go func(messageChannel chan string, wg *sync.WaitGroup) {
		for message := range messageChannel {
			file.Write([]byte(message))
		}
		wg.Done()
	}(channel, &waitWriteFile)

	logger = &loggerInfo{
		MsgChannel:    channel,
		Verbosity:     verbosity,
		File:          file,
		WaitWriteFile: &waitWriteFile,
	}

	return err
}

func Info(format string, args ...any) {
	if logger == nil {
		return
	}

	message := fmt.Sprintf(format, args...)
	logger.MsgChannel <- fmt.Sprintf(" [INFO] %s\n", message)
}

func Error(format string, args ...any) {
	if logger == nil {
		return
	}

	message := fmt.Sprintf(format, args...)
	logger.MsgChannel <- fmt.Sprintf("\033[31m [ERROR] %s\033[39m\n", message)
}

func Debug(format string, args ...any) {
	if logger == nil || logger.Verbosity == NORMAL {
		return
	}

	message := fmt.Sprintf(format, args...)
	logger.MsgChannel <- fmt.Sprintf("\033[36m [DEBUG] %s\033[39m\n", message)
}

func Close() {
	if logger == nil {
		return
	}

	close(logger.MsgChannel)
	logger.WaitWriteFile.Wait()
}
