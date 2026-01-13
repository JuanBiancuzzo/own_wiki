package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type loggerInfo struct {
	MsgChannel    chan string
	Verbosity     Verbosity
	File          *os.File
	WaitWriteFile *sync.WaitGroup
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

	channel := make(chan string, int(config.MessageCapacity))
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
		Verbosity:     config.Verbosity,
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
