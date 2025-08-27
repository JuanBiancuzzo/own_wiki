package dependencias

import (
	"fmt"
	"strings"
)

type tipoInsercion byte

const (
	TI_SELECT = iota
	TI_WHERE
)

type InformacionClave struct {
	Variable Variable
	Nombre   string
	Alias    string
	Path     []string
}

type NodoClave struct {
	Padre       *NodoClave
	Tabla       *Tabla
	Nombre      string
	Select      []HojaClave
	Where       []HojaClave
	Referencias []*NodoClave
}

type HojaClave struct {
	Padre    NodoClave
	Nombre   string
	Alias    string
	Variable Variable
}

func NewRaizClave(tabla *Tabla) NodoClave {
	return NodoClave{
		Padre:       nil,
		Tabla:       tabla,
		Nombre:      "",
		Select:      []HojaClave{},
		Where:       []HojaClave{},
		Referencias: []*NodoClave{},
	}
}

func NewNodoClave(padre *NodoClave, tabla *Tabla, nombreClave string) NodoClave {
	return NodoClave{
		Padre:       padre,
		Tabla:       tabla,
		Nombre:      nombreClave,
		Select:      []HojaClave{},
		Where:       []HojaClave{},
		Referencias: []*NodoClave{},
	}
}

func (nc NodoClave) InsertarWhere(clave string) (*HojaClave, error) {
	return nc.insertar(clave, TI_WHERE)
}

func (nc NodoClave) InsertarSelect(clave string) (*HojaClave, error) {
	return nc.insertar(clave, TI_SELECT)
}

func (nc NodoClave) insertar(clave string, tipo tipoInsercion) (*HojaClave, error) {
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

		return nodo.insertar(clave[indiceDivision:], tipo)

	} else if _, ok := variable.Informacion.(VariableArrayReferencia); ok {
		return nil, fmt.Errorf("todavia no esta soportado las array referencia")

	} else if indice, contiene := nc.ContieneClave(clave, tipo); contiene {
		switch tipo {
		case TI_SELECT:
			return &nc.Select[indice], nil
		case TI_WHERE:
			return &nc.Where[indice], nil
		default:
			return nil, fmt.Errorf("de alguna forma ")
		}

	} else {
		separacion := strings.Split(clave, "=")
		if len(separacion) > 2 {
			return nil, fmt.Errorf("se tiene para la clave %s un error de formato, donde se espera que este dado clave=alias", clave)
		}

		nodoInsertado := HojaClave{
			Nombre:   strings.TrimSpace(separacion[0]),
			Alias:    strings.TrimSpace(separacion[len(separacion)-1]),
			Variable: variable,
		}
		nc.Select = append(nc.Select, nodoInsertado)
		return &nodoInsertado, nil
	}
}

func (nc NodoClave) ContieneClave(clave string, tipo tipoInsercion) (int, bool) {
	switch tipo {
	case TI_SELECT:
		for i, claveHoja := range nc.Select {
			if claveHoja.Nombre == clave {
				return i, true
			}
		}
	case TI_WHERE:
		for i, claveHoja := range nc.Where {
			if claveHoja.Nombre == clave {
				return i, true
			}
		}
	}

	return 0, false
}

func (nc NodoClave) ObtenerPath() []string {
	if nc.Padre == nil {
		return []string{}
	}

	return append(nc.Padre.ObtenerPath(), nc.Tabla.NombreTabla)
}

func (hc HojaClave) ObtenerInfoVariable() InformacionClave {
	return InformacionClave{
		Variable: hc.Variable,
		Nombre:   hc.Nombre,
		Alias:    hc.Alias,
		Path:     hc.Padre.ObtenerPath(),
	}
}
