package tablas

import (
	"encoding/binary"
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

const LIBROS = "libros"
const TABLA_LIBROS = `CREATE TABLE libros (
  id INT AUTO_INCREMENT PRIMARY KEY,

  titulo VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  subtitulo VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  anio YEAR,
  edicion INT,
  volumen INT,
  url VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idEditorial INT,
  idArchivo INT,

  FOREIGN KEY (idEditorial) REFERENCES editoriales(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
)`
const INSERTAR_LIBRO = "INSERT INTO libros (titulo, subtitulo, anio, edicion, volumen, url) VALUES (?, ?, ?, ?, ?, ?)"

type Libros struct {
	RefEditorial *Editoriales
	RefArchivos  *Archivos
	Tracker      *d.TrackerDependencias
}

func NewLibros(refArchivos *Archivos, refEditorial *Editoriales, tracker *d.TrackerDependencias, canalMensajes chan string) (Libros, error) {
	libros := Libros{
		RefEditorial: refEditorial,
		RefArchivos:  refArchivos,
		Tracker:      tracker,
	}

	if err := tracker.RegistrarTabla(libros, d.DEPENDIENTE_DEPENDIBLE, canalMensajes); err != nil {
		return libros, fmt.Errorf("error al registrar Libros con error: %v", err)
	} else {
		return libros, nil
	}
}

func (l Libros) Nombre() string {
	return LIBROS
}

func (l Libros) CargarLibro(pathArchivo string, editorial string, titulo, subtitulo string, anio int, edicion int, volumen int, url string) error {
	hashEditorial := l.RefEditorial.GenerarHash(editorial)
	hashArchivo := l.RefArchivos.GenerarHash(pathArchivo)

	return l.Tracker.InsertarDependiente(
		l, l.GenerarHash(titulo, subtitulo, anio, edicion, volumen),
		[]d.ForeignKey{
			d.NewForeignKey("idEditorial", l.RefEditorial.Nombre(), hashEditorial),
			d.NewForeignKey("idArchivo", l.RefArchivos.Nombre(), hashArchivo),
		},
		titulo, subtitulo, anio, edicion, volumen, url,
	)
}

func (l Libros) GenerarHash(titulo, subtitulo string, anio int, edicion int, volumen int) d.IntFK {
	bufAnio := make([]byte, 4)
	binary.BigEndian.PutUint32(bufAnio, uint32(anio))

	bufEdicion := make([]byte, 4)
	binary.BigEndian.PutUint32(bufEdicion, uint32(edicion))

	bufVolumen := make([]byte, 4)
	binary.BigEndian.PutUint32(bufVolumen, uint32(volumen))

	datos := append([]byte(titulo), []byte(subtitulo)...)
	datos = append(datos, bufAnio...)
	datos = append(datos, bufEdicion...)
	datos = append(datos, bufVolumen...)

	return l.Tracker.Hash.HasearDatos(datos)
}

func (l Libros) Query(bdd *b.Bdd, datos ...any) (int64, error) {
	return bdd.Insertar(INSERTAR_LIBRO, datos...)
}

func (l Libros) ObjetoExistente(bdd *b.Bdd, datos ...any) (bool, error) {
	return false, nil
}

func (l Libros) CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error {
	if err := bdd.CrearTabla(fmt.Sprintf(TABLA_LIBROS, info.MaxNombre, info.MaxNombre, info.MaxUrl)); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de Personas, con error: %v", err)
	}
	return nil
}

func (l Libros) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (l Libros) ObtenerDependencias() []d.Tabla {
	return []d.Tabla{*l.RefArchivos, *l.RefEditorial}
}
