CREATE TABLE IF NOT EXISTS archivos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  path VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci 
);

CREATE TABLE IF NOT EXISTS personas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  apellido VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS editoriales (
  id INT AUTO_INCREMENT PRIMARY KEY,
  editorial VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS tags (
  tag VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS libros (
  id INT AUTO_INCREMENT PRIMARY KEY,
  titulo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  subtitulo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  anio YEAR,
  idEditorial INT NOT NULL,
  edicion INT,
  volumen INT,
  url VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idEditorial) REFERENCES editoriales(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS capitulos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  capitulo INT,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  paginaInicial INT,
  paginaFinal INT,
  idLibro INT NOT NULL,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idLibro) REFERENCES libros(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS autoresLibro (
  idLibro INT NOT NULL,
  idPersona INT NOT NULL,

  FOREIGN KEY (idLibro) REFERENCES libros(id),
  FOREIGN KEY (idPersona) REFERENCES personas(id)
);

CREATE TABLE IF NOT EXISTS editoresCapitulo (
  idCapitulo INT NOT NULL,
  idPersona INT NOT NULL,

  FOREIGN KEY (idCapitulo) REFERENCES capitulos(id),
  FOREIGN KEY (idPersona) REFERENCES personas(id)
);

CREATE TABLE IF NOT EXISTS distribuciones (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  tipo ENUM ("Discreta", "Continua", "Multivariada") NOT NULL,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS relacionDistribuciones (
  idDistribucionA INT NOT NULL,
  idDistribucionB INT NOT NULL,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idDistribucionA) REFERENCES distribuciones(id),
  FOREIGN KEY (idDistribucionB) REFERENCES distribuciones(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS carreras (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  etapa ENUM ("SinEmpezar", "Empezado", "Ampliar", "Terminado"),
  tieneCodigoMateria BOOLEAN,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS planesCarrera (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS cuatrimestreCarrera (
  id INT AUTO_INCREMENT PRIMARY KEY,
  anio YEAR NOT NULL,
  cuatrimestre ENUM ("Primero", "Segundo") NOT NULL
);

CREATE TABLE IF NOT EXISTS materias (
  id INT AUTO_INCREMENT PRIMARY KEY,
  idCarrera INT NOT NULL,
  idPlan INT NOT NULL,
  idCuatrimestre INT NOT NULL,
  codigo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  etapa ENUM ("SinEmpezar", "Empezado", "Ampliar", "Terminado"),
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idCarrera) REFERENCES carreras(id),
  FOREIGN KEY (idPlan) REFERENCES planesCarrera(id),
  FOREIGN KEY (idCuatrimestre) REFERENCES cuatrimestreCarrera(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS materiasEquivalentes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  codigo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idCarrera INT NOT NULL,
  idMateria INT NOT NULL,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idCarrera) REFERENCES carreras(id),
  FOREIGN KEY (idMateria) REFERENCES materias(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS materiasCorrelativas (
  tipoMateria ENUM ("Materia", "Equivalente"),
  idMateria INT NOT NULL,
  tipoCorrelativa ENUM ("Materia", "Equivalente"),
  idCorrelativa INT NOT NULL
);

CREATE TABLE IF NOT EXISTS temasMateria (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  capitulo INT,
  parte INT,
  idMateria INT NOT NULL,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idMateria) REFERENCES materias(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS paginasCursos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombrePagina VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS cursos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  etapa ENUM ("SinEmpezar", "Empezado", "Ampliar", "Terminado"),
  anioCurso YEAR,
  idPagina INT NOT NULL,
  url VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idPagina) REFERENCES paginasCursos(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS cursosPresencial (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  etapa ENUM ("SinEmpezar", "Empezado", "Ampliar", "Terminado"),
  anioCurso YEAR,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS temasCurso (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  capitulo INT,
  parte INT,
  tipoCurso ENUM ("Online", "Presencial"),
  idCurso INT NOT NULL,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS profesoresCurso (
  idCurso INT NOT NULL,
  tipoCurso ENUM ("Online", "Presencial"),
  idPersona INT NOT NULL,

  FOREIGN KEY (idPersona) REFERENCES personas(id)
);

CREATE TABLE IF NOT EXISTS temasInvestigacion (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS subtemasInvestigacion (
  idTema INT NOT NULL,
  idSubtema INT NOT NULL,

  FOREIGN KEY (idTema) REFERENCES temasInvestigacion(id),
  FOREIGN KEY (idSubtema) REFERENCES temasInvestigacion(id)
);

CREATE TABLE IF NOT EXISTS revistasDePapers (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?)  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS papers (
  id INT AUTO_INCREMENT PRIMARY KEY,
  titulo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  subtitulo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idRevista INT NOT NULL,
  volumenRevista INT,
  numeroRevista INT,
  paginaInicio INT,
  paginaFinal INT,
  anio YEAR,
  url VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idRevista) REFERENCES revistasDePapers(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS escritoresPaper (
  tipo ENUM ("Editor", "Autor"),
  idPaper INT NOT NULL,
  idPersona INT NOT NULL,

  FOREIGN KEY (idPaper) REFERENCES papers(id),
  FOREIGN KEY (idPersona) REFERENCES personas(id)
);

CREATE TABLE IF NOT EXISTS temasMatematica (
  id INT AUTO_INCREMENT PRIMARY KEY,
  numRepresentante INT,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS subtemasMatematica (
  idTema INT NOT NULL,
  idSubtema INT NOT NULL,

  FOREIGN KEY (idTema) REFERENCES temasMatematica(id),
  FOREIGN KEY (idSubtema) REFERENCES temasMatematica(id)
);

CREATE TABLE IF NOT EXISTS bloqueMatematica (
  id INT AUTO_INCREMENT PRIMARY KEY,
  idTema INT NOT NULL,
  numRepresentante INT,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  tipo ENUM ("Teorema", "Procposicion", "Observacion", "Definicion", "Colorario"),
  idArchivo INT NOT NULL,

  FOREIGN KEY (idTema) REFERENCES temasMatematica(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS colorarioBloque (
  idColorario INT NOT NULL,
  idBloque INT NOT NULL,

  FOREIGN KEY (idColorario) REFERENCES bloqueMatematica(id),
  FOREIGN KEY (idBloque) REFERENCES bloqueMatematica(id)
);

CREATE TABLE IF NOT EXISTS gruposLegales (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS seccionesLegales (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idGrupo INT NOT NULL,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idGrupo) REFERENCES gruposLegales(id),
  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS documentosLegales (
  id INT AUTO_INCREMENT PRIMARY KEY,
  abreviacion VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  articulosTienenNombre bool,
  idSeccion INT NOT NULL,

  FOREIGN KEY (idSeccion) REFERENCES seccionesLegales(id)
);

CREATE TABLE IF NOT EXISTS articulos (
  idSeccion INT NOT NULL,
  numero INT,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,

  FOREIGN KEY (idSeccion) REFERENCES seccionesLegales(id)
);

CREATE TABLE IF NOT EXISTS gruposDocumento (
  idDocumento INT NOT NULL,
  idGrupo INT NOT NULL,

  FOREIGN KEY (idDocumento) REFERENCES documentosLegales(id),
  FOREIGN KEY (idGrupo) REFERENCES gruposLegales(id)
);

CREATE TABLE IF NOT EXISTS notas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  etapa ENUM ("SinEmpezar", "Empezado", "Ampliar", "Terminado"),
  dia DATE,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
);

CREATE TABLE IF NOT EXISTS notasVinculo (
  idNota INT NOT NULL,
  idVinculo INT NOT NULL,
  tipoVinculo ENUM ("Facultad", "Coleccion", "Curso", "Investigacion", "Proyecto"),

  FOREIGN KEY (idNota) REFERENCES notas(id)
);

CREATE TABLE IF NOT EXISTS canalesYoutube (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS referenciasYoutube (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombreVideo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idCanal INT NOT NULL,
  fecha datetime,
  url VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,

  FOREIGN KEY (idCanal) REFERENCES canalesYoutube(id)
);

CREATE TABLE IF NOT EXISTS paginasWeb (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS referenciasWeb (
  id INT AUTO_INCREMENT PRIMARY KEY,
  titulo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idPagina INT NOT NULL,
  fecha datetime,
  url VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,

  FOREIGN KEY (idPagina) REFERENCES paginasWeb(id)
);

CREATE TABLE IF NOT EXISTS articulosWebAutor (
  idPaginaWeb INT NOT NULL,
  idAutor INT NOT NULL,

  FOREIGN KEY (idPaginaWeb) REFERENCES referenciasWeb(id),
  FOREIGN KEY (idAutor) REFERENCES personas(id)
);

CREATE TABLE IF NOT EXISTS referenciasWiki (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombreArticulo VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  fecha datetime,
  url VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS nombresDiccionario (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
);

CREATE TABLE IF NOT EXISTS referenciasDiccionario (
  id INT AUTO_INCREMENT PRIMARY KEY,
  palabra VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idDiccionario INT NOT NULL,
  idEditorial INT NOT NULL,
  fecha datetime,
  url VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,

  FOREIGN KEY (idDiccionario) REFERENCES nombresDiccionario(id),
  FOREIGN KEY (idEditorial) REFERENCES editoriales(id)
);

CREATE TABLE IF NOT EXISTS referencias (
  id INT AUTO_INCREMENT PRIMARY KEY,
  tipo ENUM ("Libro", "CapituloLibro", "Paper", "Curso", "TemaCurso", "Youtube", "Web", "Wiki", "Diccionario"),
  idReferencia INT NOT NULL
);