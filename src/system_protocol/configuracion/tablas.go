package configuracion

import (
	"encoding/json"
	"fmt"
	"strings"

	d "own_wiki/system_protocol/dependencias"
)

type TipoVariable string

const (
	TV_INT        = "int"
	TV_ENUM       = "enum"
	TV_STRING     = "string"
	TV_BOOL       = "bool"
	TV_DATE       = "date"
	TV_REFERENCIA = "ref"
	TV_ARRAY_REF  = "arrayRef"
)

type InfoTabla struct {
	Nombre             string             `json:"nombre"`
	Independiente      bool               `json:"independiente"`
	Dependible         bool               `json:"dependible"`
	ElementosRepetidos bool               `json:"elementosRepetidos"`
	ValoresGuardar     []InfoValorGuardar `json:"valoresGuardar"`
}

type InfoValorGuardar struct {
	Clave          string       `json:"clave"`
	Tipo           TipoVariable `json:"tipo"`
	Representativo bool         `json:"representativo"`
	Necesario      bool         `json:"necesario"`

	// Enm
	Valores []string `json:"valores"`

	// String
	Largo int `json:"largo"`

	// Referencia y ArrayReferencia
	Tabla  string             `json:"tabla"`
	Tablas []string           `json:"tablas"`
	Extra  []InfoValorGuardar `json:"valoresExtra"`
}

func CrearTablas(archivoJson string) ([]d.DescripcionTabla, error) {
	tablas := []d.DescripcionTabla{}

	decodificador := json.NewDecoder(strings.NewReader(archivoJson))

	// read open bracket
	if _, err := decodificador.Token(); err != nil {
		return tablas, err
	}

	listaInfo := []InfoTabla{}
	tipoTablas := []d.TipoTabla{}

	mapaReferenciados := make(map[string]bool)
	mapaExistencia := make(map[string]bool)

	for decodificador.More() {
		var info InfoTabla
		if err := decodificador.Decode(&info); err != nil {
			return tablas, fmt.Errorf("error al codificar tablas, con err: %v", err)
		}

		listaInfo = append(listaInfo, info)
		for _, valorGuardado := range info.ValoresGuardar {
			independiente := true
			_, dependible := mapaReferenciados[info.Nombre]

			mapaExistencia[valorGuardado.Tabla] = true

			if valorGuardado.Tipo == TV_REFERENCIA || valorGuardado.Tipo == TV_ARRAY_REF {
				independiente = false

				if valorGuardado.Tabla != "" {
					mapaReferenciados[valorGuardado.Tabla] = true
				}
				for _, tabla := range valorGuardado.Tablas {
					mapaReferenciados[tabla] = true
				}
			}

			var tipoTabla d.TipoTabla
			if independiente && dependible {
				tipoTabla = d.INDEPENDIENTE_DEPENDIBLE
			} else if independiente && !dependible {
				tipoTabla = d.INDEPENDIENTE_NO_DEPENDIBLE
			} else if !independiente && dependible {
				tipoTabla = d.DEPENDIENTE_DEPENDIBLE
			} else {
				tipoTabla = d.DEPENDIENTE_NO_DEPENDIBLE
			}
			tipoTablas = append(tipoTablas, tipoTabla)
		}
	}

	// read closing bracket
	if _, err := decodificador.Token(); err != nil {
		return tablas, err
	}

	nombresTablas := []string{}
	for nombreTabla := range mapaExistencia {
		nombresTablas = append(nombresTablas, nombreTabla)
	}

	mapaTablas := make(map[string]*d.DescripcionTabla)
	for i, info := range listaInfo {
		variables := []d.Variable{}
		for _, vg := range info.ValoresGuardar {
			necesario := vg.Necesario
			representativo := vg.Representativo && necesario
			clave := vg.Clave

			switch vg.Tipo {
			case TV_INT:
				variables = append(variables, d.NewVariableInt(representativo, clave, necesario))

			case TV_BOOL:
				variables = append(variables, d.NewVariableBool(representativo, clave, necesario))

			case TV_DATE:
				variables = append(variables, d.NewVariableDate(representativo, clave, necesario))

			case TV_STRING:
				variable := d.NewVariableString(representativo, clave, uint(vg.Largo), necesario)
				variables = append(variables, variable)

			case TV_ENUM:
				variable := d.NewVariableEnum(representativo, clave, vg.Valores, necesario)
				variables = append(variables, variable)

			case TV_REFERENCIA:
				var nombreTablas []string
				if vg.Tabla != "" {
					nombreTablas = []string{vg.Tabla}
				} else {
					nombreTablas = vg.Tablas
				}

				tablasRelacionadas := make([]*d.DescripcionTabla, len(nombreTablas))
				for i, nombreTabla := range nombreTablas {
					if tabla, ok := mapaTablas[nombreTabla]; !ok {
						return tablas, fmt.Errorf("la tabla %s no esta registrada, esto puede ser un error de tipeo, ya que el resto de las tablas son: [%s]", nombreTabla, strings.Join(nombresTablas, ", "))
					} else {
						tablasRelacionadas[i] = tabla
					}
				}

				variables = append(variables, d.NewVariableReferencia(vg.Representativo, clave, tablasRelacionadas))

			case TV_ARRAY_REF:
				var nombreTablas []string
				if vg.Tabla != "" {
					nombreTablas = []string{vg.Tabla}
				} else {
					nombreTablas = vg.Tablas
				}

				tablasRelacionadas := make([]*d.DescripcionTabla, len(nombreTablas))
				for i, nombreTabla := range nombreTablas {
					if nombreTabla == info.Nombre {
						tablasRelacionadas[i] = nil
						continue
					}

					if tabla, ok := mapaTablas[nombreTabla]; !ok {
						return tablas, fmt.Errorf("la tabla %s no esta registrada, esto puede ser un error de tipeo, ya que el resto de las tablas son: [%s]", nombreTabla, strings.Join(nombresTablas, ", "))
					} else {
						tablasRelacionadas[i] = tabla
					}
				}

				variables = append(variables, d.NewVariableArrayReferencias(clave, tablasRelacionadas))

			default:
				return tablas, fmt.Errorf("el tipo de dato %s no existe, debe ser un error", vg.Tipo)
			}

		}

		nuevaTabla := d.ConstruirTabla(info.Nombre, tipoTablas[i], info.ElementosRepetidos, variables)
		mapaTablas[info.Nombre] = &nuevaTabla

		tablas = append(tablas, nuevaTabla)
	}

	return tablas, nil
}
