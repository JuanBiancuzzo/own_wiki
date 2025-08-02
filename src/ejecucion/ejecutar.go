package ejecucion

import (
	"bufio"
	"fmt"
	"os"
	"own_wiki/ejecucion/fs"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

func MostrarLs(directorio *fs.Directorio, canalMensajes chan string) {
	if output, err := directorio.Ls(); err != nil {
		canalMensajes <- fmt.Sprintf("Se intentó ejecutar ls y se obtuvo un error: %v", err)
	} else {
		canalMensajes <- output
	}
}

func EjectuarCd(directorio *fs.Directorio, parametroCd string, canalMensajes chan string) {
	if err := directorio.Cd(parametroCd); err != nil {
		canalMensajes <- fmt.Sprintf("Se intentó ejecutar cd con '%s' y se obtuvo un error: %v", parametroCd, err)
	}
}

func Abs(valor int) int {
	if valor < 0 {
		return -valor
	}
	return valor
}

func Ejecutar(dirInput string, canalMensajes chan string) {
	directorio, err := fs.NewDirectorio()
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer directorio.Close()

	scanner := bufio.NewScanner(os.Stdin)
	canalMensajes <- "Terminal prendida..."

	for {
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if Abs(strings.Index(input, "ls")) == 0 {
			MostrarLs(directorio, canalMensajes)

		} else if Abs(strings.Index(input, "cd ")) == 0 {
			EjectuarCd(directorio, strings.TrimSpace(input[3:]), canalMensajes)

		} else if input == "cd" {
			EjectuarCd(directorio, "../../../../..", canalMensajes)

		} else {
			canalMensajes <- "El comando '%s' no puede no es ls o cd"
		}

	}
}
