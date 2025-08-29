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

/*
	"Carreras": {
		"tipo": "ElementoUnico",
		"tabla": "Carreras",
		"id": "idCarrera",
		"claves": [ "nombre" ]
	}

Resulta en la query

	"SELECT nombre FROM Carreras WHERE id = ?"

Para esto, usamos la variable de "nombre", y se tiene que buscar el indice del parametro "idCarrera" para
pasarlo en el lugar de ?

Otro ejemplo:

	"Materia": {
		"tipo": "ElementoUnico",
		"tabla": "Materias",
		"id": "idMateria",
		"claves": [ "id", "nombre", "refCarrera:id", "refCarrera:nombre" ]
	},

Resulta en la query

	"SELECT res.Materias_id, res.Materias_nombre, res.Carreras_id, res.Carreras_nombre FROM (
		SELECT Materias.id AS Materias_id, Materias.nombre AS Materias_nombre, Carreras.id AS Carreras_id, Carreras.nombre AS Carreras_nombre
		FROM Materias INNER JOIN Carreras ON Materias.refCarrera == Carreras.id
	) AS res WHERE res.Materias_id = ?"

Otro ejemplo:

	"Tema": {
	    "tipo": "ElementoUnico",
	    "tabla": "TemasMateria",
	    "id": "idTema",
	    "claves": [ "nombre", "refMateria:id", "refMateria:nombre", "refMateria:refCarrera:id", "refMateria:refCarrera:nombre" ]
	},

# Resulta en la query

	Todas las claves serian:
		[ TemasMateria_id, TemasMateria_nombre, temp1_1.Materias_id, temp_1_1.Materias_nombre, temp_2_1.Carreras_id, temp_2_1.Carreras_nombre ]
	por lo tanto se puede hacer un wrapper de esta sentencia, con todos los

	SELECT TemasMateria.id AS TemasMateria_id, TemasMateria.nombre AS TemasMateria_nombre, temp1_1.* FROM TemasMateria
	JOIN (
		SELECT Materias.id AS Materias_id, Materias.nombre AS Materias_nombre, temp2_1.* FROM Materias
		JOIN (
			SELECT Carreras.id AS Carreras_id, Carreras.nombre AS Carreras_nombre FROM Carreras
		)
		AS temp2_1 ON Materias.refCarrera = temp2_1.Carreras_id
	)
	AS temp1_1 ON TemasMateria.refMateria = temp1_1.Materias_id
	WHERE TemasMateria_id = ?;
*/
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

/*
	type InformacionQuery struct {
	    Condicion    string
	    OrderBy      []string
	    ClavesUsadas []string
	}

	"Materias": {
	    "Materias": {
	        "condicion": "refCarrera:id == idCarrera",
	        "orderBy": [ "refCuatrimestre:anio=anio", "refCuatrimestre:cuatrimestre=cuatrimestre" ],
	        "claves": [ "id", "nombre", "refCarrera:id", "refCuatrimestre:anio", "refCuatrimestre:cuatrimestre" ]
	    },
	    "MateriasEquivalente": {
	        "condicion": "refCarrera:id == idCarrera",
	        "orderBy": [ "refMateria:refCuatrimestre:anio=anio", "refMateria:refCuatrimestre:cuatrimestre=cuatrimestre" ],
	        "claves": [ "nombre", "refCarrera:id", "refMateria:id", "refMateria:refCuatrimestre:anio", "refMateria:refCuatrimestre:cuatrimestre" ]
	    }
	}
*/
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
