package ejecucion

import (
	"fmt"
	fs "own_wiki/ejecucion/fs"

	_ "github.com/go-sql-driver/mysql"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

func MostrarLs(directorio *fs.Directorio, canalMensajes chan string) bool {
	if output, err := directorio.Ls(); err != nil {
		canalMensajes <- fmt.Sprintf("Se intentó ejecutar ls y se obtuvo un error: %v", err)
		return false

	} else {
		canalMensajes <- output
		return true
	}
}

func EjectuarCd(directorio *fs.Directorio, parametroCd string, canalMensajes chan string) bool {
	if err := directorio.Cd(parametroCd); err != nil {
		canalMensajes <- fmt.Sprintf("Se intentó ejecutar cd con '%s' y se obtuvo un error: %v", parametroCd, err)
		return false
	}
	return true
}

func Ejecutar(dirInput string, canalMensajes chan string) {
	directorio, err := fs.NewDirectorio()
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer directorio.Close()

	if !MostrarLs(directorio, canalMensajes) {
		return
	}

	canalMensajes <- ""
	canalMensajes <- ""
	canalMensajes <- ""

	if !EjectuarCd(directorio, fs.PD_FACULTAD, canalMensajes) {
		return
	}

	if !MostrarLs(directorio, canalMensajes) {
		return
	}
}
