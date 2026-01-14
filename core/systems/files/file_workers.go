package files

import (
	"sync"

	log "github.com/JuanBiancuzzo/own_wiki/core/systems/logger"
)

type FnWorkFile func(File)

const DEFAULT_FILE_AMOUNT = 20

func newFileWorker(workFiles FnWorkFile, waitFiles *sync.WaitGroup) chan string {
	filePaths := make(chan string, DEFAULT_FILE_AMOUNT)

	waitFiles.Add(1)
	go func(filePaths chan string, workFiles FnWorkFile, waitFiles *sync.WaitGroup) {
		for filePath := range filePaths {
			if file, err := NewFile(filePath); err != nil {
				log.Warn("Failed to create file info: %q, with error: %v. Skipping file", filePath, err)

			} else {
				workFiles(file)
			}
		}

		waitFiles.Done()
	}(filePaths, workFiles, waitFiles)

	return filePaths
}

func WorkFilesRoundRobin(filePaths chan string, amountWorkers uint, workFile FnWorkFile, waitFiles *sync.WaitGroup) {
	channelsFilePaths := make([]chan string, amountWorkers)
	for i := range amountWorkers {
		channelsFilePaths[i] = newFileWorker(workFile, waitFiles)
	}

	var workerToSend uint = 0
	for filePath := range filePaths {
		channelsFilePaths[workerToSend] <- filePath
		workerToSend = (workerToSend + 1) % amountWorkers
	}

	for i := range amountWorkers {
		close(channelsFilePaths[i])
	}
}
