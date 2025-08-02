package estructura

import (
	"database/sql"
	"fmt"
	"strconv"
)

const INSERTAR_PAPER = "INSERT INTO papers (titulo, subtitulo, idRevista, volumenRevista, numeroRevista, paginaInicio, paginaFinal, anio, url, idArchivo) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

const QUERY_REVISTA_PAPER = "SELECT id FROM revistasDePapers WHERE nombre = ?"
const INSERTAR_REVISTA_PAPER = "INSERT INTO revistasDePapers (nombre) VALUES (?)"

const INSERTAR_ESCRITOR_PAPER = "INSERT INTO escritoresPaper (tipo, idPaper, idPersona) VALUES (?, ?, ?)"

type TipoEscritorPaper string

const (
	PAPER_EDITOR = "Editor"
	PAPER_AUTOR  = "Autor"
)

type Paper struct {
	IdArchivo      *Opcional[int64]
	Titulo         string
	Subtitulo      string
	NombreRevista  string
	VolumenRevista int
	NumeroRevista  int
	PaginaInicio   int
	PaginaFinal    int
	Anio           int
	Url            string
	Autores        []*Persona
	Editores       []*Persona
}

func NewPaper(titulo string, subtitulo string, nombreRevista string, volumenRevista string, numeroRevista string, paginaInicial string, paginaFinal string, repAnio string, url string, autores []*Persona, editores []*Persona) (*Paper, error) {
	if anio, err := strconv.Atoi(repAnio); err != nil {
		return nil, fmt.Errorf("error al crear paper al obtener el anio, con error: %v", err)
	} else {
		return &Paper{
			IdArchivo:      NewOpcional[int64](),
			Titulo:         titulo,
			Subtitulo:      subtitulo,
			NombreRevista:  nombreRevista,
			VolumenRevista: NumeroODefault(volumenRevista, 0),
			NumeroRevista:  NumeroODefault(numeroRevista, 0),
			PaginaInicio:   NumeroODefault(paginaInicial, 0),
			PaginaFinal:    NumeroODefault(paginaFinal, 0),
			Anio:           anio,
			Url:            url,
			Autores:        autores,
			Editores:       editores,
		}, nil
	}
}

func (p *Paper) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		p.IdArchivo.Asignar(id)
		return p, true
	})
}

func (p *Paper) Insertar(idRevista int64) ([]any, error) {
	if idArchivo, existe := p.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("paper no tiene todavia el idArchivo")

	} else {
		return []any{p.Titulo, p.Subtitulo, idRevista, p.VolumenRevista, p.NumeroRevista, p.PaginaInicio, p.PaginaFinal, p.Anio, p.Url, idArchivo}, nil
	}
}

func (p *Paper) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- "Insertar Paper"

	var idPaper int64
	if idRevista, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_REVISTA_PAPER, p.NombreRevista) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_REVISTA_PAPER, p.NombreRevista) },
	); err != nil {
		return 0, fmt.Errorf("error al hacer una querry de la revista %s con error: %v", p.NombreRevista, err)

	} else if datos, err := p.Insertar(idRevista); err != nil {
		return 0, err

	} else if idPaper, err = InsertarDirecto(bdd, INSERTAR_PAPER, datos...); err != nil {
		return 0, err
	}

	for _, escritor := range p.Autores {
		if idAutor, err := ObtenerOInsertar(
			func() *sql.Row { return bdd.QueryRow(QUERY_PERSONAS, escritor.Insertar()...) },
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, escritor.Insertar()...) },
		); err != nil {
			canal <- fmt.Sprintf("error al hacer una querry del autor: %s %s con error: %v", escritor.Nombre, escritor.Apellido, err)

		} else if _, err := bdd.Exec(INSERTAR_ESCRITOR_PAPER, PAPER_AUTOR, idPaper, idAutor); err != nil {
			canal <- fmt.Sprintf("error al insertar par idRevista-idEscritor, con error: %v", err)
		}
	}

	for _, escritor := range p.Editores {
		if idAutor, err := ObtenerOInsertar(
			func() *sql.Row { return bdd.QueryRow(QUERY_PERSONAS, escritor.Insertar()...) },
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, escritor.Insertar()...) },
		); err != nil {
			canal <- fmt.Sprintf("error al hacer una querry del autor: %s %s con error: %v", escritor.Nombre, escritor.Apellido, err)

		} else if _, err := bdd.Exec(INSERTAR_ESCRITOR_PAPER, PAPER_EDITOR, idPaper, idAutor); err != nil {
			canal <- fmt.Sprintf("error al insertar par idRevista-idEscritor, con error: %v", err)
		}
	}

	return idPaper, nil
}

func (p *Paper) ResolverDependencias(id int64) []Cargable {
	return []Cargable{}
}

func (p *Paper) CargarDependencia(dependencia Dependencia) {}
