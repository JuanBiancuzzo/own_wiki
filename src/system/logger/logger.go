package logger

import (
	"fmt"
	"os"
	"sync"
)

type Verbosity string

const (
	MUTE    = "mute"
	NORMAL  = "normal"
	VERBOSE = "verbose"
)

type loggerInfo struct {
	Channel   chan string
	Verbosity Verbosity
	File      *os.File
	WaitGroup *sync.WaitGroup
}

var logger *loggerInfo = nil

func CreateLogger(path string, verbosity Verbosity) (err error) {
	if verbosity == MUTE {
		return nil
	}

	var file *os.File = nil
	if file, err = os.Create(path); err != nil {
		return fmt.Errorf("failed to open logger file: %w", err)
	}

	channel := make(chan string)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func(messageChannel chan string, wg *sync.WaitGroup) {
		for message := range messageChannel {
			file.Write([]byte(message))
		}
		wg.Done()
	}(channel, &waitGroup)

	logger = &loggerInfo{
		Channel:   channel,
		Verbosity: verbosity,
		File:      file,
		WaitGroup: &waitGroup,
	}

	return err
}

func Info(string string, args ...any) {
	if logger == nil {
		return
	}

	message := fmt.Sprintf(string, args...)
	logger.Channel <- fmt.Sprintf(" [INFO] %s\n", message)
}

func Error(string string, args ...any) {
	if logger == nil {
		return
	}

	message := fmt.Sprintf(string, args...)
	logger.Channel <- fmt.Sprintf("\033[31m [ERROR] %s\033[39m\n", message)
}

func Debug(string string, args ...any) {
	if logger == nil || logger.Verbosity == NORMAL {
		return
	}

	message := fmt.Sprintf(string, args...)
	logger.Channel <- fmt.Sprintf("\033[36m [DEBUG] %s\033[39m\n", message)
}

func Close() {
	if logger == nil {
		return
	}

	close(logger.Channel)
	logger.WaitGroup.Wait()
}
