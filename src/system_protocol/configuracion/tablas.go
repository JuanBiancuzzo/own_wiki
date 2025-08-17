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

	// Referencia
	Tabla  string   `json:"tabla"`
	Tablas []string `json:"tablas"`
}

func CrearTablas(archivoJson string) ([]d.DescripcionTabla, error) {
	tablas := []d.DescripcionTabla{}

	decodificador := json.NewDecoder(strings.NewReader(archivoJson))

	// read open bracket
	if _, err := decodificador.Token(); err != nil {
		return tablas, err
	}

	listaInfo := []InfoTabla{}
	mapaReferenciados := make(map[string]uint8)

	for decodificador.More() {
		var info InfoTabla
		if err := decodificador.Decode(&info); err != nil {
			return tablas, fmt.Errorf("error al codificar tablas, con err: %v", err)
		}

		listaInfo = append(listaInfo, info)
		for _, valorGuardado := range info.ValoresGuardar {
			if valorGuardado.Tipo != TV_REFERENCIA {
				continue
			}

			if valorGuardado.Tabla != "" {
				mapaReferenciados[valorGuardado.Tabla] = 0
			}
			for _, tabla := range valorGuardado.Tablas {
				mapaReferenciados[tabla] = 0
			}
		}
	}

	// read closing bracket
	if _, err := decodificador.Token(); err != nil {
		return tablas, err
	}

	mapaTablas := make(map[string]*d.DescripcionTabla)
	for _, info := range listaInfo {
		independiente := true
		_, dependible := mapaReferenciados[info.Nombre]

		paresClaveTipo := []d.ParClaveTipo{}
		referenciasTablas := []d.ReferenciaTabla{}
		for _, vg := range info.ValoresGuardar {
			var nuevoClaveTipo d.ParClaveTipo

			necesario := vg.Necesario
			representativo := vg.Representativo && necesario

			switch vg.Tipo {
			case TV_STRING:
				nuevoClaveTipo = d.NewClaveString(representativo, vg.Clave, uint(vg.Largo), necesario)
				paresClaveTipo = append(paresClaveTipo, nuevoClaveTipo)

			case TV_INT:
				nuevoClaveTipo = d.NewClaveInt(representativo, vg.Clave, necesario)
				paresClaveTipo = append(paresClaveTipo, nuevoClaveTipo)

			case TV_ENUM:
				nuevoClaveTipo = d.NewClaveEnum(representativo, vg.Clave, vg.Valores, necesario)
				paresClaveTipo = append(paresClaveTipo, nuevoClaveTipo)

			case TV_BOOL:
				nuevoClaveTipo = d.NewClaveBool(representativo, vg.Clave, necesario)
				paresClaveTipo = append(paresClaveTipo, nuevoClaveTipo)

			case TV_DATE:
				nuevoClaveTipo = d.NewClaveDate(representativo, vg.Clave, necesario)
				paresClaveTipo = append(paresClaveTipo, nuevoClaveTipo)

			case TV_REFERENCIA:
				independiente = false
				var nombreTablas []string
				if vg.Tabla != "" {
					nombreTablas = []string{vg.Tabla}
				} else {
					nombreTablas = vg.Tablas
				}

				tablasRelacionadas := make([]*d.DescripcionTabla, len(nombreTablas))
				for i, nombreTabla := range nombreTablas {
					if tabla, ok := mapaTablas[nombreTabla]; !ok {
						nombreTablas := []string{}
						for nombreTabla := range mapaTablas {
							nombreTablas = append(nombreTablas, nombreTabla)
						}
						return tablas, fmt.Errorf("la tabla %s no esta registrada, esto puede ser un error de tipeo, ya que el resto de las tablas son: [%s]", nombreTabla, strings.Join(nombreTablas, ", "))
					} else {
						tablasRelacionadas[i] = tabla
					}
				}
				nuevaReferencia := d.NewReferenciaTabla(vg.Clave, tablasRelacionadas, vg.Representativo)
				referenciasTablas = append(referenciasTablas, nuevaReferencia)

			default:
				return tablas, fmt.Errorf("el tipo de dato %s no existe, debe ser un error", vg.Tipo)
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

		nuevaTabla := d.ConstruirTabla(info.Nombre, tipoTabla, info.ElementosRepetidos, paresClaveTipo, referenciasTablas)
		mapaTablas[info.Nombre] = &nuevaTabla

		tablas = append(tablas, nuevaTabla)
	}

	return tablas, nil
}
