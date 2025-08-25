package dependencias

import (
	"fmt"
	"strings"
)

// "tabla": "TemasMateria",
// "claves": [ "nombre", "refMateria:id", "refMateria:nombre", "refMateria:refCarrera:id", "refMateria:refCarrera:nombre" ]

type InformacionClave struct {
	Variable TipoVariable
	Nombre   string
	Path     []string
}

type NodoClave struct {
	Padre       *NodoClave
	Tabla       *DescripcionTabla
	Nombre      string
	Claves      []HojaClave
	Referencias []*NodoClave
}

func NewRaizClave(tabla *DescripcionTabla) NodoClave {
	return NodoClave{
		Padre:       nil,
		Tabla:       tabla,
		Nombre:      "",
		Claves:      []HojaClave{},
		Referencias: []*NodoClave{},
	}
}

func NewNodoClave(padre *NodoClave, tabla *DescripcionTabla, nombreClave string) NodoClave {
	return NodoClave{
		Padre:       padre,
		Tabla:       tabla,
		Nombre:      nombreClave,
		Claves:      []HojaClave{},
		Referencias: []*NodoClave{},
	}
}

func (nc NodoClave) Insertar(clave string) (*HojaClave, error) {
	indiceDivision := strings.Index(clave, ":")
	primeraClave := clave
	if indiceDivision > 0 {
		primeraClave = clave[:indiceDivision]
	}
	variable, ok := nc.Tabla.Variables[primeraClave]
	if !ok {
		return nil, fmt.Errorf("la clave %s no existe en la tabla %s", clave, nc.Tabla.NombreTabla)
	}

	if info, ok := variable.Informacion.(VariableReferencia); ok {
		if len(info.Tablas) > 1 {
			return nil, fmt.Errorf("todavia no se puede referenciar multiples tablas")
		}
		tabla := info.Tablas[0]

		var nodo *NodoClave = nil
		for _, referencia := range nc.Referencias {
			if referencia.Tabla.NombreTabla == primeraClave {
				nodo = referencia
				break
			}
		}

		if nodo == nil {
			nuevoNodo := NewNodoClave(&nc, tabla, primeraClave)
			nc.Referencias = append(nc.Referencias, &nuevoNodo)
			nodo = &nuevoNodo
		}

		return nodo.Insertar(clave[indiceDivision:])

	} else if _, ok := variable.Informacion.(VariableArrayReferencia); ok {
		return nil, fmt.Errorf("todavia no esta soportado las array referencia")

	} else if tipo, err := variable.ObtenerTipo(); err != nil {
		return nil, fmt.Errorf("al obtener el tipo se tuvo el error: %v", err)

	} else {
		nodoInsertado := HojaClave{
			Nombre:   clave,
			Variable: tipo,
		}
		nc.Claves = append(nc.Claves, nodoInsertado)
		return &nodoInsertado, nil
	}
}

func (nc NodoClave) ObtenerPath() []string {
	if nc.Padre == nil {
		return []string{}
	}

	return append(nc.Padre.ObtenerPath(), nc.Tabla.NombreTabla)
}

type HojaClave struct {
	Padre    NodoClave
	Nombre   string
	Variable TipoVariable
}

func (hc HojaClave) ObtenerInfoVariable() InformacionClave {
	return InformacionClave{
		Variable: hc.Variable,
		Nombre:   hc.Nombre,
		Path:     hc.Padre.ObtenerPath(),
	}
}

func (hc HojaClave) NombreQuery() string {
	return fmt.Sprintf("%s_%s", hc.Padre.Tabla.NombreTabla, hc.Nombre)
}
