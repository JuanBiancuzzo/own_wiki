package estructura

import (
	"database/sql"
	"fmt"
	"strings"

	l "own_wiki/system_protocol/listas"
)

const INSERTAR_ARCHIVO = "INSERT INTO archivos (path) VALUES (?)"
const QUERY_ARCHIVO = "SELECT id FROM archivos WHERE path = ?"

const INSERTAR_TAG = "INSERT INTO tags (tag, idArchivo) VALUES (?, ?)"

type Archivo struct {
	Path              string
	Tags              []string
	ListaDependencias *l.Lista[Dependencia]
}

func NewArchivo(path string, tags []string) *Archivo {
	return &Archivo{
		Path:              path,
		Tags:              tags,
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (a *Archivo) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	if idArchivo, err := Insertar(func() (sql.Result, error) { return bdd.Exec(INSERTAR_ARCHIVO, a.Path) }); err != nil {
		return 0, fmt.Errorf("error al obtener insertar el archivo: %s, con error: %v", Nombre(a.Path), err)

	} else {
		for _, tag := range a.Tags {
			if _, err := bdd.Exec(INSERTAR_TAG, tag, idArchivo); err != nil {
				canal <- fmt.Sprintf("Error al insertar tag: %s en el archivo: %s\n", tag, Nombre(a.Path))
			}
		}

		return idArchivo, nil
	}
}

func (a *Archivo) ResolverDependencias(id int64) []Cargable {
	cantidadCumple := 0
	cargables := make([]Cargable, a.ListaDependencias.Largo)

	for cumpleDependencia := range a.ListaDependencias.Iterar {
		if cargable, cumple := cumpleDependencia(id); cumple {
			cargables[cantidadCumple] = cargable
			cantidadCumple++
		}
	}

	return cargables[:cantidadCumple]
}

func (a *Archivo) CargarDependencia(dependencia Dependencia) {
	a.ListaDependencias.Push(dependencia)
}

func Nombre(path string) string {
	separacion := strings.Split(path, "/")
	return strings.ReplaceAll(separacion[len(separacion)-1], ".md", "")
}
