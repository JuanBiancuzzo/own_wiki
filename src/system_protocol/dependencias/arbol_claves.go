package dependencias

import (
	"fmt"
	"strings"
)

type tipoInsercion string

const (
	TI_SELECT = "select"
	TI_WHERE  = "where"
)

type InformacionClave struct {
	Variable DescripcionVariable
	Nombre   string
	Alias    string
	Path     []string
}

type NodoClave struct {
	Padre       *NodoClave
	Tabla       *DescripcionTabla
	Nombre      string
	Select      []HojaClave
	Where       []HojaClave
	Referencias []*NodoClave
}

type HojaClave struct {
	Padre    *NodoClave
	Nombre   string
	Alias    string
	Variable DescripcionVariable
}

func NewRaizClave(tabla *DescripcionTabla) *NodoClave {
	return &NodoClave{
		Padre:       nil,
		Tabla:       tabla,
		Nombre:      "",
		Select:      []HojaClave{},
		Where:       []HojaClave{},
		Referencias: []*NodoClave{},
	}
}

func NewNodoClave(padre *NodoClave, tabla *DescripcionTabla, nombreClave string) *NodoClave {
	return &NodoClave{
		Padre:       padre,
		Tabla:       tabla,
		Nombre:      nombreClave,
		Select:      []HojaClave{},
		Where:       []HojaClave{},
		Referencias: []*NodoClave{},
	}
}

func (nc *NodoClave) InsertarWhere(clave string, tablas map[string]*DescripcionTabla) (*HojaClave, error) {
	// fmt.Printf("Insertando where clave: '%s'\n", strings.TrimSpace(clave))
	return nc.insertar(strings.TrimSpace(clave), TI_WHERE, tablas)
}

func (nc *NodoClave) InsertarSelect(clave string, tablas map[string]*DescripcionTabla) (*HojaClave, error) {
	// fmt.Printf("Insertando select clave: '%s'\n", strings.TrimSpace(clave))
	return nc.insertar(strings.TrimSpace(clave), TI_SELECT, tablas)
}

func (nc *NodoClave) insertar(clave string, tipo tipoInsercion, tablas map[string]*DescripcionTabla) (*HojaClave, error) {
	indiceDivision := strings.Index(clave, ":")
	primeraClave := clave
	resto := ""
	if indiceDivision > 0 {
		primeraClave = strings.TrimSpace(clave[:indiceDivision])
		resto = strings.TrimSpace(clave[indiceDivision+1:])
	}

	if strings.Contains(primeraClave, ".") {
		separacion := strings.Split(primeraClave, ".")
		primeraClave = separacion[0]
		resto = separacion[1]
	}

	// fmt.Printf("Viendo la subclave: '%s'\n", primeraClave)

	variable, ok := nc.Tabla.ObtenerVariable(primeraClave)
	if !ok {
		return nil, fmt.Errorf("la clave %s no existe en la tabla %s", clave, nc.Tabla.Nombre)
	}

	if info, ok := variable.Descripcion.(DescVariableReferencia); ok {
		// fmt.Printf("Procesando nodo: '%s'\n", primeraClave)
		if len(info.Tablas) > 1 {
			return nil, fmt.Errorf("todavia no se puede referenciar multiples tablas")
		}

		var nodo *NodoClave = nil
		for _, referencia := range nc.Referencias {
			if referencia.Nombre == primeraClave {
				// fmt.Println("Ya tenemos ese nodo")
				nodo = referencia
				break
			}
		}

		// usando unicamente el primero por ahora
		if tabla, ok := tablas[info.Tablas[0]]; ok && nodo == nil {
			// fmt.Println("No lo tenemos, vamos a crear nodo")
			nodo = NewNodoClave(nc, tabla, primeraClave)
			nc.Referencias = append(nc.Referencias, nodo)
		}

		if indiceDivision < 0 {
			resto = "id"
		}

		// fmt.Printf("Insertando: '%s'\n", claveInsertar)
		return nodo.insertar(strings.TrimSpace(resto), tipo, tablas)

	} else if _, ok := variable.Descripcion.(DescVariableArrayReferencia); ok {
		return nil, fmt.Errorf("todavia no esta soportado las array referencia")

	} else if indice, contiene := nc.ContieneClave(clave, tipo); contiene {
		// fmt.Printf("Ya se insertó (%s) la clave: '%s'\n", tipo, clave)
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
		// fmt.Printf("Insertando el nodo: '%s' con alias: '%s'\n", separacion[0], separacion[len(separacion)-1])

		nodoInsertado := HojaClave{
			Padre:    nc,
			Nombre:   strings.TrimSpace(separacion[0]),
			Alias:    strings.TrimSpace(separacion[len(separacion)-1]),
			Variable: variable,
		}
		switch tipo {
		case TI_SELECT:
			nc.Select = append(nc.Select, nodoInsertado)
		case TI_WHERE:
			nc.Where = append(nc.Where, nodoInsertado)
		}
		return &nodoInsertado, nil
	}
}

func (nc NodoClave) ObtenerClaveSelect() []*HojaClave {
	hojas := make([]*HojaClave, len(nc.Select))
	for i, hoja := range nc.Select {
		hojas[i] = &hoja
	}

	for _, referencia := range nc.Referencias {
		hojas = append(hojas, referencia.ObtenerClaveSelect()...)
	}

	return hojas
}

func (nc NodoClave) ObtenerClaveWhere() []*HojaClave {
	hojas := make([]*HojaClave, len(nc.Where))
	for i, hoja := range nc.Where {
		hojas[i] = &hoja
	}

	for _, referencia := range nc.Referencias {
		hojas = append(hojas, referencia.ObtenerClaveWhere()...)
	}

	return hojas
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

	return append(nc.Padre.ObtenerPath(), nc.Nombre)
}

func (hc HojaClave) ObtenerInfoVariable() InformacionClave {
	return InformacionClave{
		Variable: hc.Variable,
		Nombre:   hc.Nombre,
		Alias:    hc.Alias,
		Path:     hc.Padre.ObtenerPath(),
	}
}
