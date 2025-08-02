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

type ConstructorPaper struct {
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

func NewConstructorPaper(titulo string, subtitulo string, nombreRevista string, volumenRevista string, numeroRevista string, paginaInicial string, paginaFinal string, repAnio string, url string, autores []*Persona, editores []*Persona) (*ConstructorPaper, error) {
	if anio, err := strconv.Atoi(repAnio); err != nil {
		return nil, fmt.Errorf("error al crear paper al obtener el anio, con error: %v", err)
	} else {
		return &ConstructorPaper{
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

func (cp *ConstructorPaper) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		return &Paper{
			Titulo:         cp.Titulo,
			Subtitulo:      cp.Subtitulo,
			NombreRevista:  cp.NombreRevista,
			VolumenRevista: cp.VolumenRevista,
			NumeroRevista:  cp.NumeroRevista,
			PaginaInicio:   cp.PaginaInicio,
			PaginaFinal:    cp.PaginaFinal,
			Anio:           cp.Anio,
			Url:            cp.Url,
			Autores:        cp.Autores,
			Editores:       cp.Editores,
			IdArchivo:      id,
		}, true
	})
}

type Paper struct {
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
	IdArchivo      int64
}

func (p *Paper) Insertar(idRevista int64) []any {
	return []any{p.Titulo, p.Subtitulo, idRevista, p.VolumenRevista, p.NumeroRevista, p.PaginaInicio, p.PaginaFinal, p.Anio, p.Url, p.IdArchivo}
}

func (p *Paper) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- "Insertar Paper"

	var idPaper int64

	if idRevista, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_REVISTA_PAPER, p.NombreRevista) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_REVISTA_PAPER, p.NombreRevista) },
	); err != nil {
		return 0, fmt.Errorf("error al hacer una querry de la revista %s con error: %v", p.NombreRevista, err)

	} else if idPaper, err = Insertar(func() (sql.Result, error) {
		return bdd.Exec(INSERTAR_PAPER, p.Insertar(idRevista)...)
	}); err != nil {
		return 0, fmt.Errorf("error al insertar un paper, con error: %v", err)
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
