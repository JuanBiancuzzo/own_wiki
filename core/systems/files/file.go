package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Name      string
	Directory string
	Content   string
	ModTime   time.Time
}

func NewFile(filePath string) (File, error) {
	directory, name := filepath.Split(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return File{}, fmt.Errorf("Failed to open file, with error: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return File{}, fmt.Errorf("Failed to get file stats, with error: %v", err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return File{}, fmt.Errorf("Failed to read the file content, with error: %v", err)
	}

	return File{
		Name:      name,
		Directory: directory,
		Content:   string(content),
		ModTime:   fileInfo.ModTime(),
	}, nil
}
