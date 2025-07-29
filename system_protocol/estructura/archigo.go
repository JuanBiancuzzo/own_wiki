package estructura

import (
	"database/sql"
	"fmt"
	"strings"
)

const INSERTAR_ARCHIVO = "INSERT INTO archivos (path) VALUES (?)"
const QUERY_ARCHIVO = "SELECT id FROM archivos WHERE path = ?"

const INSERTAR_TAG = "INSERT INTO tags (tag, idArchivo) VALUES (?, ?)"

type Archivo struct {
	Path string
	Tags []string
}

func NewArchivo(path string, tags []string) *Archivo {
	return &Archivo{
		Path: path,
		Tags: tags,
	}
}

func (a *Archivo) CargarDatos(bdd *sql.DB, canal chan string) bool {
	if idArchivo, err := Insertar(func() (sql.Result, error) { return bdd.Exec(INSERTAR_ARCHIVO, a.Path) }); err != nil {
		canal <- fmt.Sprintf("Error al obtener insertar el archivo: %s, con error: %v\n", Nombre(a.Path), err)
	} else {
		for _, tag := range a.Tags {
			if _, err := bdd.Exec(INSERTAR_TAG, tag, idArchivo); err != nil {
				canal <- fmt.Sprintf("Error al insertar tag: %s en el archivo: %s\n", tag, Nombre(a.Path))
			}
		}
	}

	return true
}

func Nombre(path string) string {
	separacion := strings.Split(path, "/")
	return strings.ReplaceAll(separacion[len(separacion)-1], ".md", "")
}
