package file_loader

import "sync"

/*
Vamos a definir una gorutina que este leyendo todos los archivos en un directorio dado, mandandolo
a un canal los file path.

Definimos otra funcion que tenga workers que dado el path de un file, lo lea y obtenga toda la
metadata del archivo, y el string del archivo para pasarlo al usuario
*/

func GetFileFromDirectory(dirPath string, files chan string) error {
	return nil
}

const MESSAGE_CAPACITY = 10

type File struct {
	Data string
	// Metadata
}

func NewFile(filePath string) File {
	return File{}
}

type ProcessFile func(file File)

func NewReaderWorker(amount int, filePaths chan string, process ProcessFile, wg *sync.WaitGroup) {
	waitGroups := make([]*sync.WaitGroup, amount)
	workers := make([]chan string, amount)

	for i := range amount {
		worker := make(chan string, MESSAGE_CAPACITY)
		var waitGroup sync.WaitGroup

		waitGroup.Add(1)
		go func(worker chan string, process ProcessFile, wg *sync.WaitGroup) {
			for filePath := range worker {
				process(NewFile(filePath))
			}
			wg.Done()
		}(worker, process, &waitGroup)

		workers[i] = worker
		waitGroups[i] = &waitGroup
	}

	wg.Add(1)
	go func(filePaths chan string, wg *sync.WaitGroup) {
		// Ejecutamos la positica de round robin
		position := 0
		for filePath := range filePaths {
			workers[position] <- filePath
			position = (position + 1) % amount
		}

		// Para este punto se cierra el filePaths y por lo tanto se toma que termino de limpiarse
		for _, worker := range workers {
			close(worker)
		}

		for _, waitGroup := range waitGroups {
			waitGroup.Wait()
		}

		wg.Done()
	}(filePaths, wg)
}
