package fs

import (
	"database/sql"
	"fmt"
	bdd "own_wiki/system_protocol/baseDeDatos"
	"strings"
)

type Directorio struct {
	Subpath      Subpath
	Bdd          *sql.DB
	CacheSubpath *Cache
}

type Subpath interface {
	Ls() ([]string, error)
	Cd(subpath string, cache *Cache) (Subpath, error)
}

func NewDirectorio(canalMensajes chan string) (*Directorio, error) {
	if bdd, err := bdd.EstablecerConexionRelacional(canalMensajes); err != nil {
		return nil, fmt.Errorf("no se pudo establecer la conexion con la base de datos, con error: %v", err)

	} else {
		cache := NewCache(bdd)
		if root, err := cache.ObtenerSubpath(PD_ROOT); err != nil {
			return nil, err
		} else {

			return &Directorio{
				Subpath:      root,
				Bdd:          bdd,
				CacheSubpath: cache,
			}, nil
		}
	}
}

func (d *Directorio) Ls() (string, error) {
	resultado := "./\n../\n"
	if lineas, err := d.Subpath.Ls(); err != nil {
		return resultado, err
	} else {
		for _, linea := range lineas {
			resultado += strings.TrimSpace(linea) + "\n"
		}
	}

	return resultado, nil
}

func (d *Directorio) Cd(path string) error {
	for subpath := range strings.SplitSeq(path, "/") {
		if subpath == "." {
			continue
		}

		var err error
		if d.Subpath, err = d.Subpath.Cd(subpath, d.CacheSubpath); err != nil {
			return err
		}
	}

	return nil
}

func (d *Directorio) Close() {
	d.Bdd.Close()
}
