package configuracion

import (
	"encoding/json"
	"fmt"
	"strings"

	d "github.com/JuanBiancuzzo/own_wiki/system_protocol/dependencias"
	u "github.com/JuanBiancuzzo/own_wiki/system_protocol/utilidades"
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

	// Array de referencias
	Estructura []InfoValorGuardar `json:"estructura"`
}

func procesarInformacionTabla(info InfoTabla, colaTablas *u.Cola[InfoTabla]) (d.DescripcionTabla, error) {

	variables := []d.DescripcionVariable{}
	for _, vg := range info.ValoresGuardar {
		necesario := vg.Necesario
		representativo := vg.Representativo && necesario
		clave := vg.Clave

		switch vg.Tipo {
		case TV_INT:
			variables = append(variables, d.NewDescVariableSimple(d.TVS_INT, representativo, clave, necesario))

		case TV_BOOL:
			variables = append(variables, d.NewDescVariableSimple(d.TVS_BOOL, representativo, clave, necesario))

		case TV_DATE:
			variables = append(variables, d.NewDescVariableSimple(d.TVS_DATE, representativo, clave, necesario))

		case TV_STRING:
			variable := d.NewDescVariableString(representativo, clave, uint(vg.Largo), necesario)
			variables = append(variables, variable)

		case TV_ENUM:
			variable := d.NewDescVariableEnum(representativo, clave, vg.Valores, necesario)
			variables = append(variables, variable)

		case TV_REFERENCIA:
			var nombreTablas []string
			if strings.TrimSpace(vg.Tabla) != "" {
				nombreTablas = []string{vg.Tabla}

			} else {
				nombreTablas = make([]string, len(vg.Tablas))
				copy(nombreTablas, vg.Tablas)
			}
			variables = append(variables, d.NewDescVariableReferencia(vg.Representativo, clave, nombreTablas))

		case TV_ARRAY_REF:
			tablaReferenciada := fmt.Sprintf("array%s_%s", clave, info.Nombre)
			selfClave := "selfRef"

			colaTablas.Encolar(InfoTabla{
				Nombre:             tablaReferenciada,
				ElementosRepetidos: false,
				ValoresGuardar: append([]InfoValorGuardar{
					{
						Clave:          selfClave,
						Tipo:           TV_REFERENCIA,
						Representativo: true,
						Necesario:      true,
						Tabla:          info.Nombre,
					},
				}, vg.Estructura...),
			})

			variables = append(variables, d.NewDescVariableArrayReferencias(clave, selfClave, tablaReferenciada))

		default:
			return d.DescripcionTabla{}, fmt.Errorf("el tipo de dato %s no existe, debe ser un error", vg.Tipo)
		}
	}

	return d.NewDescripcionTabla(
		info.Nombre,
		info.ElementosRepetidos,
		variables,
	), nil
}

func DescribirTablas(archivoJson string) ([]d.DescripcionTabla, error) {
	descripcionTablas := []d.DescripcionTabla{}
	decodificador := json.NewDecoder(strings.NewReader(archivoJson))

	// read open bracket
	if _, err := decodificador.Token(); err != nil {
		return descripcionTablas, err
	}

	extraTablas := u.NewCola[InfoTabla]()

	for decodificador.More() {
		var info InfoTabla
		if err := decodificador.Decode(&info); err != nil {
			return descripcionTablas, fmt.Errorf("error al codificar tablas, con err: %v", err)
		}

		if descripcion, err := procesarInformacionTabla(info, extraTablas); err != nil {
			return descripcionTablas, fmt.Errorf("tuvo un error al procesar, con error: %v", err)

		} else {
			descripcionTablas = append(descripcionTablas, descripcion)
		}
	}

	// read closing bracket
	if _, err := decodificador.Token(); err != nil {
		return descripcionTablas, err
	}

	for !extraTablas.Lista.Vacia() {
		nuevasTablasProcesar := extraTablas.Lista.Items()
		extraTablas.Lista.Vaciar()

		for _, info := range nuevasTablasProcesar {
			if descripcion, err := procesarInformacionTabla(info, extraTablas); err != nil {
				return descripcionTablas, fmt.Errorf("tuvo un error al procesar, con error: %v", err)

			} else {
				descripcionTablas = append(descripcionTablas, descripcion)
			}
		}
	}

	return descripcionTablas, nil
}

func CrearTablas(archivoJson string, tracker *d.TrackerDependencias) ([]d.Tabla, error) {
	tablas := []d.Tabla{}

	descripcionTablas, err := DescribirTablas(archivoJson)
	if err != nil {
		return tablas, err
	}

	mapaReferenciados := make(map[string]bool)
	for _, descTabla := range descripcionTablas {
		for _, descVariable := range descTabla.Variables {
			if detalle, ok := descVariable.Descripcion.(d.DescVariableReferencia); ok {
				for _, tabla := range detalle.Tablas {
					mapaReferenciados[tabla] = true
				}
			}
		}
	}

	mapaTablas := make(map[string]*d.Tabla)
	for _, descTabla := range descripcionTablas {
		independiente := true
		_, dependible := mapaReferenciados[descTabla.Nombre]

		variables := make([]d.Variable, len(descTabla.Variables))
		for i, descVariable := range descTabla.Variables {

			switch detalle := descVariable.Descripcion.(type) {
			case d.DescVariableSimple:
				variables[i] = d.NewVariableSimple(detalle.Tipo, detalle.Representativo, descVariable.Clave, detalle.Necesario)

			case d.DescVariableString:
				variables[i] = d.NewVariableString(detalle.Representativo, descVariable.Clave, detalle.Largo, detalle.Necesario)

			case d.DescVariableEnum:
				variables[i] = d.NewVariableEnum(detalle.Representativo, descVariable.Clave, detalle.Valores, detalle.Necesario)

			case d.DescVariableReferencia:
				independiente = false

				tablasRelacionadas := make([]*d.Tabla, len(detalle.Tablas))
				for i, nombreTabla := range detalle.Tablas {
					if tabla, ok := mapaTablas[nombreTabla]; !ok {
						return tablas, fmt.Errorf("la tabla %s no esta registrada, esto puede ser un error de tipeo", nombreTabla)
					} else {
						tablasRelacionadas[i] = tabla
					}
				}

				variables[i] = d.NewVariableReferencia(detalle.Representativo, descVariable.Clave, tablasRelacionadas)

			case d.DescVariableArrayReferencia:
				variables[i] = d.NewVariableArrayReferencias(descVariable.Clave, detalle.ClaveSelf, detalle.TablaCreada)
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

		if nuevaTabla, err := d.ConstruirTabla(tracker, descTabla.Nombre, tipoTabla, descTabla.ElementosRepetidos, variables); err != nil {
			return tablas, err

		} else {
			mapaTablas[descTabla.Nombre] = &nuevaTabla
			tablas = append(tablas, nuevaTabla)
		}
	}

	return tablas, nil
}
