package files

import (
	"fmt"
	"os"
	"path/filepath"

	q "github.com/JuanBiancuzzo/own_wiki/core/systems/lists"
	log "github.com/JuanBiancuzzo/own_wiki/core/systems/logger"
)

type FnIgnoreFile func(filePath string) (toBeIgnore bool)

func DefaultIgnoreFileFunction() FnIgnoreFile {
	return func(string) bool { return false }
}

type FnIgnoreDirectory func(directoy string) (toBeIgnore bool)

func DefaultIgnoreDirectoryFunction() FnIgnoreDirectory {
	return func(string) bool { return false }
}

func FilesInDirectory(directoryOrigin string, filePaths chan string, ignoreFiles FnIgnoreFile, ignoreDirectories FnIgnoreFile) error {
	directoryQueue := q.NewQueue[string]()
	directoryQueue.Queue(directoryOrigin)

	log.Debug("Start searching though the directory: %s", directoryOrigin)
	defer close(filePaths)

	for directoryPath := range directoryQueue.Iterate {
		if ignoreDirectories(directoryPath) {
			continue
		}

		directory, err := os.ReadDir(directoryPath)
		if err != nil {
			return fmt.Errorf("Falied to read directory %q, with error: %v", directoryPath, err)
		}

		for _, file := range directory {
			filePath := filepath.Join(directoryPath, file.Name())
			if file.IsDir() && !ignoreDirectories(filePath) {
				directoryQueue.Queue(filePath)

			} else if !file.IsDir() && !ignoreFiles(filePath) {
				filePaths <- filePath
			}
		}
	}

	return nil
}
