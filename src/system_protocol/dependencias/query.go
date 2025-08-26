package dependencias

import (
	"fmt"
	"slices"
	"strings"
)

type QueryDato struct {
	Select string
	Claves []*HojaClave
	Where  []*HojaClave
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

	SELECT TemasMateria.id AS TemasMateria_id, TemasMateria.nombre AS TemasMateria_nombre, temp1_1.* FROM TemasMateria
	JOIN (
		SELECT Materias.id AS Materias_id, Materias.nombre AS Materias_nombre, temp2_1.* FROM Materias
		JOIN (
			SELECT Carreras.id AS Carreras_id, Carreras.nombre AS Carreras_nombre FROM Carreras
		)
		AS temp2_1 ON Materias.refCarrera = temp2_1.Carreras_id
		JOIN (
			SELECT Carreras.id AS Carreras_id, Carreras.nombre AS Carreras_nombre FROM Carreras
		)
		AS temp2_2 ON Materias.refCarrera = temp2_1.Carreras_id
	)
	AS temp1_1 ON TemasMateria.refMateria = temp1_1.Materias_id
	WHERE TemasMateria_id = 1;
*/
func generarSelect(nodo *NodoClave, profundidad int) string {
	nombreTabla := nodo.Tabla.NombreTabla
	claves := make([]string, len(nodo.Claves))
	for i, clave := range nodo.Claves {
		nombreClave := clave.Nombre
		claves[i] = fmt.Sprintf("%s.%s AS %s_%s", nombreTabla, nombreClave, nombreTabla, nombreClave)
	}

	if len(nodo.Referencias) == 0 {
		return fmt.Sprintf("SELECT %s FROM %s", strings.Join(claves, ", "), nombreTabla)
	}

	sentenciasJoin := make([]string, len(nodo.Referencias))
	for i, referencia := range nodo.Referencias {
		sentenciaInterna := generarSelect(referencia, profundidad+1)

		sentenciasJoin[i] = fmt.Sprintf(
			"INNER JOIN (\n\t%s\n) AS temp_%d_%d ON %s.%s  temp_%d_%d.%s_id",
			sentenciaInterna,
			profundidad, i,
			nodo.Tabla.NombreTabla,
			referencia.Nombre,
			profundidad, i,
			referencia.Tabla.NombreTabla,
		)
	}

	return fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(claves, ", "),
		nombreTabla,
		strings.Join(sentenciasJoin, "\n"),
	)
}

func NewQuerySimple(tabla *DescripcionTabla, clavesUsadas []string) (QueryDato, error) {
	indiceId := slices.Index(clavesUsadas, "id")
	if indiceId < 0 {
		indiceId = len(clavesUsadas)
		clavesUsadas = append(clavesUsadas, "id")
	}

	claves := make([]*HojaClave, len(clavesUsadas))

	raiz := NewRaizClave(tabla)
	for i, clave := range clavesUsadas {
		if hoja, err := raiz.Insertar(clave); err != nil {
			return QueryDato{}, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)

		} else {
			claves[i] = hoja
		}
	}

	return QueryDato{
		Select: fmt.Sprintf(
			"%s WHERE %s = ?",
			generarSelect(&raiz, 0),
			claves[indiceId].NombreQuery(),
		),
		Claves: claves,
		Where:  []*HojaClave{},
	}, nil
}

func NewQueryMultiplesCompleto(tabla *DescripcionTabla, clavesUsadas []string, condicion string) (QueryDato, error) {
	return QueryDato{}, nil
}

func NewQueryMultiplesParcial() (QueryDato, error) {
	return QueryDato{}, nil
}
