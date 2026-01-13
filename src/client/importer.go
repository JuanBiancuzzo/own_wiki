package main

import "sync"

type Importer struct {
	FilePaths chan string
	WaitGroup *sync.WaitGroup
}

func NewImporter() *Importer {
	var waitGroup sync.WaitGroup
	filePaths := make(chan string)

	return &Importer{
		FilePaths: filePaths,
		WaitGroup: &waitGroup,
	}
}

func (i *Importer) Close() {
	close(i.FilePaths)
	i.WaitGroup.Wait()
}
