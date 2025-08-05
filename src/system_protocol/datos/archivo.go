package datos

import (
	"fmt"
	"strings"

	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
)

const INSERTAR_ARCHIVO = "INSERT INTO archivos (path) VALUES (?)"
const QUERY_ARCHIVO = "SELECT id FROM archivos WHERE path = ?"

const INSERTAR_TAG = "INSERT INTO tags (tag, idArchivo) VALUES (?, ?)"

type Archivo struct {
	Path              string
	Tags              []string
	ListaDependencias *u.Lista[Dependencia]
}

func NewArchivo(path string, tags []string) *Archivo {
	return &Archivo{
		Path:              path,
		Tags:              tags,
		ListaDependencias: u.NewLista[Dependencia](),
	}
}

func (a *Archivo) CargarDatos(bdd *b.Bdd, canal chan string) (int64, error) {
	if idArchivo, err := InsertarDirecto(bdd, INSERTAR_ARCHIVO, a.Path); err != nil {
		return 0, fmt.Errorf("error al obtener insertar el archivo: %s, con error: %v", Nombre(a.Path), err)

	} else {
		for _, tag := range a.Tags {
			if _, err := InsertarDirecto(bdd, INSERTAR_TAG, tag, idArchivo); err != nil {
				canal <- fmt.Sprintf("Error al insertar tag: %s en el archivo: %s\n", tag, Nombre(a.Path))
			}
		}

		return idArchivo, nil
	}
}

func (a *Archivo) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, a.ListaDependencias.Items())
}

func (a *Archivo) CargarDependencia(dependencia Dependencia) {
	a.ListaDependencias.Push(dependencia)
}

func Nombre(path string) string {
	separacion := strings.Split(path, "/")
	return strings.ReplaceAll(separacion[len(separacion)-1], ".md", "")
}
