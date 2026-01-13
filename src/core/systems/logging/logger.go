package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Verbosity uint

const (
	MUTE = iota
	NORMAL
	VERBOSE
)

func (v Verbosity) String() string {
	switch v {
	case MUTE:
		return "Mute"

	case NORMAL:
		return "Normal"

	case VERBOSE:
		return "Verbose"

	default:
		return fmt.Sprintf("[ERROR] '%d' is not a posible Verbosity", v)
	}
}

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

	splitPath := strings.Split(path, "/")
	folder := strings.Join(splitPath[:len(splitPath)-1], "/")

	if err = os.MkdirAll(folder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to open/create logger folder, with error: %w", err)
	}

	var file *os.File = nil
	if file, err = os.Create(path); err != nil {
		return fmt.Errorf("failed to open logger file, with error: %w", err)
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
