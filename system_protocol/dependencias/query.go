package dependencias

import (
	"fmt"
	"strings"
)

type QueryDato struct {
	SentenciaQuery string
	ClaveSelect    []*HojaClave
	ClaveWhere     []*HojaClave
	Parametros     []string
}

type InformacionQuery struct {
	Condiciones  []string // Claves de la tabla
	Parametros   []string // valores pasados
	OrderBy      []string
	ClavesUsadas []string
}

func generarSetencia(nodo *NodoClave, profundidad int) string {
	nombreTabla := nodo.Tabla.Nombre
	clavesSelect := make([]string, len(nodo.Select))
	for i, clave := range nodo.Select {
		nombreClave := clave.Nombre
		alias := fmt.Sprintf("%s_%s", nombreTabla, nombreClave)
		if clave.Nombre != clave.Alias {
			alias = clave.Alias
		}

		clavesSelect[i] = fmt.Sprintf("%s.%s AS %s", nombreTabla, nombreClave, alias)
	}

	sentenciaSelect := strings.Join(clavesSelect, ", ")
	if len(nodo.Select) == 0 {
		sentenciaSelect = "*"
	}

	clavesWhere := make([]string, len(nodo.Where))
	for i, clave := range nodo.Where {
		clavesWhere[i] = fmt.Sprintf("%s.%s = ?", nombreTabla, clave.Nombre)
	}

	sentenciaWhere := ""
	if len(clavesWhere) > 0 {
		sentenciaWhere = fmt.Sprintf("WHERE %s", strings.Join(clavesWhere, " AND "))
	}

	if len(nodo.Referencias) == 0 {
		return fmt.Sprintf("SELECT %s FROM %s %s", sentenciaSelect, nombreTabla, sentenciaWhere)
	}

	sentenciasJoin := make([]string, len(nodo.Referencias))
	nombresTemporales := make([]string, len(nodo.Referencias))

	for i, referencia := range nodo.Referencias {
		sentenciaInterna := generarSetencia(referencia, profundidad+1)

		nombreTemporal := fmt.Sprintf("temp_%d_%d", profundidad, i)
		claveReferencia := fmt.Sprintf("%s.%s", nodo.Tabla.Nombre, referencia.Nombre)
		claveId := fmt.Sprintf("%s.%s_id", nombreTemporal, referencia.Tabla.Nombre)

		sentenciasJoin[i] = fmt.Sprintf(
			"INNER JOIN (\n\t%s\n) AS %s ON %s = %s",
			sentenciaInterna, nombreTemporal, claveReferencia, claveId,
		)
		nombresTemporales[i] = fmt.Sprintf("%s.*", nombreTemporal)
	}

	return fmt.Sprintf(
		"SELECT %s, %s FROM %s %s %s",
		sentenciaSelect,
		strings.Join(nombresTemporales, ", "),
		nombreTabla,
		strings.Join(sentenciasJoin, "\n"),
		sentenciaWhere,
	)
}

func NewQuerySimple(tabla *DescripcionTabla, clavesUsadas []string, parametroId string, descripciones map[string]*DescripcionTabla) (QueryDato, error) {
	var err error

	raiz := NewRaizClave(tabla)
	for _, clave := range clavesUsadas {
		if _, err = raiz.InsertarSelect(clave, descripciones); err != nil {
			return QueryDato{}, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)
		}
	}

	if _, err = raiz.InsertarSelect("id", descripciones); err != nil {
		return QueryDato{}, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)
	}

	if _, err = raiz.InsertarWhere("id", descripciones); err != nil {
		return QueryDato{}, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)
	}

	return QueryDato{
		SentenciaQuery: generarSetencia(raiz, 0),
		ClaveSelect:    raiz.ObtenerClaveSelect(),
		ClaveWhere:     raiz.ObtenerClaveWhere(),
		Parametros:     []string{parametroId},
	}, nil
}

func NewQueryMultiples(tablas map[*DescripcionTabla]InformacionQuery, groupBy []string, descripciones map[string]*DescripcionTabla) (map[string]QueryDato, error) {
	datosQuery := make(map[string]QueryDato)

	for tabla := range tablas {
		info := tablas[tabla]
		var err error

		raiz := NewRaizClave(tabla)
		for _, clave := range info.ClavesUsadas {
			if _, err = raiz.InsertarSelect(clave, descripciones); err != nil {
				return datosQuery, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)
			}
		}

		for _, clave := range info.Condiciones {
			if _, err = raiz.InsertarSelect(clave, descripciones); err != nil {
				return datosQuery, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)
			}

			if _, err = raiz.InsertarWhere(clave, descripciones); err != nil {
				return datosQuery, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)
			}
		}

		datosQuery[tabla.Nombre] = QueryDato{
			SentenciaQuery: generarSetencia(raiz, 0),
			ClaveSelect:    raiz.ObtenerClaveSelect(),
			ClaveWhere:     raiz.ObtenerClaveWhere(),
			Parametros:     info.Parametros,
		}
	}

	return datosQuery, nil
}
