package encoding

import (
	"fmt"
	"strings"

	fs "own_wiki/encoding/fs"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"

	_ "embed"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"

//go:embed tablas.json
var infoTablas string

type InfoTabla struct {
	Nombre             string                `json:"nombre"`
	Independiente      bool                  `json:"independiente"`
	Dependible         bool                  `json:"dependible"`
	ElementosRepetidos bool                  `json:"elementosRepetidos"`
	ValoresGuardar     []InfoValorGuardar    `json:"valoresGuardar"`
	ReferenciasTabla   []InfoReferenciaTabla `json:"referenciasTabla"`
}

type InfoValorGuardar struct {
	Clave          string   `json:"clave"`
	Tipo           string   `json:"tipo"`
	Largo          int      `json:"largo"`
	Valores        []string `json:"valores"`
	Representativo bool     `json:"representativo"`
	Necesario      bool     `json:"necesario"`
}

type InfoReferenciaTabla struct {
	Tabla          string   `json:"tabla"`
	Tablas         []string `json:"tablas"`
	Clave          string   `json:"clave"`
	Representativo bool     `json:"representativo"`
}

func CrearTablas() ([]d.DescripcionTabla, error) {
	tablas := []d.DescripcionTabla{}

	decodificador := json.NewDecoder(strings.NewReader(infoTablas))

	// read open bracket
	if _, err := decodificador.Token(); err != nil {
		return tablas, err
	}

	listaInfo := []InfoTabla{}
	mapaReferenciados := make(map[string]uint8)

	for decodificador.More() {
		var info InfoTabla
		err := decodificador.Decode(&info)
		if err != nil {
			return tablas, fmt.Errorf("error al codificar tablas")
		}

		listaInfo = append(listaInfo, info)
		for _, referencia := range info.ReferenciasTabla {
			mapaReferenciados[referencia.Tabla] = 0
			for _, tabla := range referencia.Tablas {
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
		independiente := len(info.ReferenciasTabla) == 0
		_, dependible := mapaReferenciados[info.Nombre]

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

		paresClaveTipo := []d.ParClaveTipo{}
		for _, vg := range info.ValoresGuardar {
			var nuevoClaveTipo d.ParClaveTipo

			necesario := vg.Necesario
			representativo := vg.Representativo && necesario

			switch vg.Tipo {
			case "string":
				nuevoClaveTipo = d.NewClaveString(representativo, vg.Clave, uint(vg.Largo), necesario)

			case "int":
				nuevoClaveTipo = d.NewClaveInt(representativo, vg.Clave, necesario)

			case "enum":
				nuevoClaveTipo = d.NewClaveEnum(representativo, vg.Clave, vg.Valores, necesario)

			case "bool":
				nuevoClaveTipo = d.NewClaveBool(representativo, vg.Clave, necesario)

			case "date":
				nuevoClaveTipo = d.NewClaveDate(representativo, vg.Clave, necesario)

			default:
				return tablas, fmt.Errorf("el tipo de dato %s no existe, debe ser un error", vg.Tipo)
			}

			paresClaveTipo = append(paresClaveTipo, nuevoClaveTipo)
		}

		referenciasTablas := []d.ReferenciaTabla{}
		for _, rt := range info.ReferenciasTabla {
			var nombreTablas []string
			if rt.Tabla != "" {
				nombreTablas = []string{rt.Tabla}
			} else {
				nombreTablas = rt.Tablas
			}

			tablasRelacionadas := make([]*d.DescripcionTabla, len(nombreTablas))
			for i, nombreTabla := range nombreTablas {
				if tabla, ok := mapaTablas[nombreTabla]; !ok {
					nombreTablas := []string{}
					for nombreTabla := range mapaTablas {
						nombreTablas = append(nombreTablas, nombreTabla)
					}
					return tablas, fmt.Errorf("la tabla %s no esta registrada, esto puede ser un error de tipeo, ya que el resto de las tablas son: [%s]", rt.Tabla, strings.Join(nombreTablas, ", "))
				} else {
					tablasRelacionadas[i] = tabla
				}
			}
			nuevaReferencia := d.NewReferenciaTabla(rt.Clave, tablasRelacionadas, rt.Representativo)
			referenciasTablas = append(referenciasTablas, nuevaReferencia)
		}

		nuevaTabla := d.ConstruirTabla(info.Nombre, tipoTabla, info.ElementosRepetidos, paresClaveTipo, referenciasTablas)
		mapaTablas[info.Nombre] = &nuevaTabla

		tablas = append(tablas, nuevaTabla)
	}

	return tablas, nil
}

func Encodear(dirInput string, canalMensajes chan string) {
	_ = godotenv.Load()

	bddRelacional, err := b.EstablecerConexionRelacional(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer b.CerrarBddRelacional(bddRelacional)

	bddNoSQL, err := b.EstablecerConexionNoSQL(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer b.CerrarBddNoSQL(bddNoSQL)
	canalMensajes <- "Se conectaron correctamente las bdd necesarias"

	tablas, err := CrearTablas()
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear las tablas, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se leyeron correctamente las tablas"

	tracker, err := d.NewTrackerDependencias(b.NewBdd(bddRelacional, bddNoSQL), tablas, canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear el tracker, se tuvo el error: %v", err)
		return
	}

	if err = fs.RecorrerDirectorio(dirInput, tracker, canalMensajes); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo recorrer todos los archivos, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se termino el proceso de insertar datos"

	if err = tracker.TerminarProcesoInsertarDatos(); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo terminar el proceso de insertar datos, se tuvo el error: %v", err)
	} else {
		canalMensajes <- "Se termino de cargar a la base de datos"
	}
	canalMensajes <- "Se hizo la limpieza de los datos auxiliares"
}
