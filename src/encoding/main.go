package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"own_wiki/encoding/procesar"
	b "own_wiki/system_protocol/base_de_datos"
	c "own_wiki/system_protocol/configuracion"
	d "own_wiki/system_protocol/dependencias"
	u "own_wiki/system_protocol/utilidades"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"

const CANTIDAD_WORKERS = 20

var DIRECTORIOS_IGNORAR = []string{".git", ".configuracion", ".github", ".obsidian", ".trash"}

func ContarArchivos(dirOrigen string) (int, error) {
	colaDirectorios := u.NewCola[string]()
	colaDirectorios.Encolar("")
	cantidadArchivos := 0

	for directorioPath := range colaDirectorios.DesencolarIterativamente {
		archivos, err := os.ReadDir(fmt.Sprintf("%s/%s", dirOrigen, directorioPath))
		if err != nil {
			return 0, fmt.Errorf("se tuvo un error al leer el directorio '%s' dando el error: %v", fmt.Sprintf("%s/%s", dirOrigen, directorioPath), err)
		}

		for _, archivo := range archivos {
			nombreArchivo := archivo.Name()
			archivoPath := nombreArchivo
			if directorioPath != "" {
				archivoPath = fmt.Sprintf("%s/%s", directorioPath, archivoPath)
			}

			if archivo.IsDir() && !slices.Contains(DIRECTORIOS_IGNORAR, nombreArchivo) {
				colaDirectorios.Encolar(archivoPath)

			} else if !archivo.IsDir() {
				cantidadArchivos++
			}
		}
	}

	return cantidadArchivos, nil
}

func MensajeProcesados(procesados, cantidadArchivos int) string {
	porcentaje := int((100 * procesados) / cantidadArchivos)
	return fmt.Sprintf(
		"Progreso: [%s%s] %d%s %04d/%04d",
		strings.Repeat("|", porcentaje), strings.Repeat(" ", 100-porcentaje),
		porcentaje, "%",
		procesados, cantidadArchivos,
	)
}

func ContadorArchivos(dirOrigen string, canalSacar, canalDone chan bool, canalMensajes chan string) {
	type contador struct {
		cantidad int
		err      error
	}

	canalCantidad := make(chan contador)
	go func(dirOrigen string, canalCantidad chan contador) {
		cantidadArchivos, err := ContarArchivos(dirOrigen)
		canalCantidad <- contador{
			cantidad: cantidadArchivos,
			err:      err,
		}
	}(dirOrigen, canalCantidad)

	seguir := true
	procesados := 0
	cantidadArchivos := 0

	for seguir {
		select {
		case <-canalSacar:
			procesados++

		case contador := <-canalCantidad:
			if contador.err != nil {
				canalMensajes <- fmt.Sprintf("No se puede hacer el conteo, se obtuvo el error: %v", contador.err)
				return
			}
			cantidadArchivos = contador.cantidad
			seguir = false

		case <-canalDone:
			canalMensajes <- MensajeProcesados(procesados, procesados)
			return
		}
	}

	canalMensajes <- MensajeProcesados(procesados, cantidadArchivos)

	seguir = true
	porcentajePrevio := 0
	for seguir {
		select {
		case <-canalSacar:
			procesados++
			porcentaje := int((100 * procesados) / cantidadArchivos)

			if procesados%10 == 0 {
				canalMensajes <- MensajeProcesados(procesados, cantidadArchivos)
			}

			if porcentaje != porcentajePrevio {
				porcentajePrevio = porcentaje
			}

		case <-canalDone:
			seguir = false
		}
	}

	canalMensajes <- MensajeProcesados(procesados, cantidadArchivos)
}

func RecorrerDirectorio(dirOrigen string, tracker *d.TrackerDependencias, canalMensajes chan string) error {
	var waitArchivos sync.WaitGroup
	canalMensajes <- fmt.Sprintf("El directorio para trabajar va a ser: %s\n", dirOrigen)

	canalInput := make(chan string, CANTIDAD_WORKERS)
	waitArchivos.Add(1)

	canalSacar := make(chan bool, 20*CANTIDAD_WORKERS)
	canalDone := make(chan bool)
	go ContadorArchivos(dirOrigen, canalSacar, canalDone, canalMensajes)

	terminar := func(motivo string) {
		canalMensajes <- "Esperando a que termine por: " + motivo
		close(canalInput)
		waitArchivos.Wait()
		close(canalSacar)
		canalDone <- true
		close(canalDone)
	}

	procesarArchivo := func(path string) {
		if err := procesar.CargarArchivo(dirOrigen, path, tracker, canalMensajes); err != nil {
			canalMensajes <- fmt.Sprintf("Se tuvo un error al crear un archivo en el path: '%s', con error: %v", path, err)
		}
		canalSacar <- true
	}
	go u.DividirTrabajo(canalInput, CANTIDAD_WORKERS, procesarArchivo, &waitArchivos)

	colaDirectorios := u.NewCola[string]()
	colaDirectorios.Encolar("")

	for directorioPath := range colaDirectorios.DesencolarIterativamente {
		archivos, err := os.ReadDir(fmt.Sprintf("%s/%s", dirOrigen, directorioPath))
		if err != nil {
			terminar("huvo un error")
			return fmt.Errorf("se tuvo un error al leer el directorio '%s' dando el error: %v", fmt.Sprintf("%s/%s", dirOrigen, directorioPath), err)
		}

		for _, archivo := range archivos {
			nombreArchivo := archivo.Name()
			archivoPath := nombreArchivo
			if directorioPath != "" {
				archivoPath = fmt.Sprintf("%s/%s", directorioPath, archivoPath)
			}

			if archivo.IsDir() && !slices.Contains(DIRECTORIOS_IGNORAR, nombreArchivo) {
				colaDirectorios.Encolar(archivoPath)

			} else if !archivo.IsDir() {
				canalInput <- archivoPath
			}
		}
	}

	terminar("termino como lo esperado")
	return nil
}

func CargarTablas(dirConfiguracion string, tracker *d.TrackerDependencias) error {
	if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", dirConfiguracion, "tablas.json")); err != nil {
		return fmt.Errorf("error al leer el archivo de configuracion para las tablas, con error: %v", err)

	} else if tablas, err := c.CrearTablas(string(bytes), tracker); err != nil {
		return err

	} else {
		for _, tabla := range tablas {
			tracker.CargarTabla(tabla)
		}
		return nil
	}
}

func Encodear(dirInput, dirOutput, dirConfiguracion string, canalMensajes chan string) {
	_ = godotenv.Load()

	bdd, err := b.NewBdd(dirOutput, canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer bdd.Close()

	tracker, err := d.NewTrackerDependencias(bdd)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear el tracker, se tuvo el error: %v", err)
		return
	}

	if err = CargarTablas(dirConfiguracion, tracker); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear las tablas, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se leyeron correctamente las tablas"

	if err = tracker.EmpezarProcesoInsertarDatos(canalMensajes); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo iniciar el proceso de insercion de datos, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Iniciando el proceso de insertar datos"

	if err = RecorrerDirectorio(dirInput, tracker, canalMensajes); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo recorrer todos los archivos, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se termino el proceso de insertar datos"

	if err = tracker.TerminarProcesoInsertarDatos(); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo terminar el proceso de insertar datos, se tuvo el error: %v", err)
	} else {
		canalMensajes <- "Se termino de cargar a la base de datos"
	}
	canalMensajes <- "Se hizo la limpieza de los datos auxiliares"
}

func main() {
	tiempoInicial := time.Now()
	var waitMensajes sync.WaitGroup
	canalMensajes := make(chan string, 100)

	waitMensajes.Add(1)
	go func(canal chan string, wg *sync.WaitGroup) {
		for mensaje := range canal {
			fmt.Println(strings.TrimSpace(mensaje))
		}
		wg.Done()
	}(canalMensajes, &waitMensajes)

	var carpetaDatos string
	var carpetaConfiguracion string
	var carpetaOutput string

	argumentoProcesar := 1
	for argumentoProcesar+1 < len(os.Args) {
		switch os.Args[argumentoProcesar] {
		case "-d":
			argumentoProcesar++
			carpetaDatos = os.Args[argumentoProcesar]
		case "-c":
			argumentoProcesar++
			carpetaConfiguracion = os.Args[argumentoProcesar]
		case "-o":
			argumentoProcesar++
			carpetaOutput = os.Args[argumentoProcesar]
		default:
			canalMensajes <- fmt.Sprintf("el argumento %s no pudo ser identificado", os.Args[argumentoProcesar])
		}
		argumentoProcesar++
	}

	configuracionValida := true
	if carpetaDatos == "" {
		canalMensajes <- "Necesitas pasar el directorio de datos (con la flag -d)"
		configuracionValida = false
	}

	if carpetaConfiguracion == "" {
		canalMensajes <- "Necesitas pasar el directorio de configuracion (con la flag -c)"
		configuracionValida = false
	}

	if carpetaOutput == "" {
		canalMensajes <- "Necesitas pasar el directorio de output (con la flag -o)"
		configuracionValida = false
	}

	if configuracionValida {
		Encodear(carpetaDatos, carpetaOutput, carpetaConfiguracion, canalMensajes)
	}

	close(canalMensajes)
	waitMensajes.Wait()

	fmt.Printf("Se termino el programa en: %s \n", time.Since(tiempoInicial))
}
