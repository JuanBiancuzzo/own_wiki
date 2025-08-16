package contenido

// TODO: Tiene celdas, y titulos, y columnas => ver que onda
/*
	Contenedores:
	Table => tal vez ver una mejor solucion
	 * TableBody
	 * TableCell
	 * Tablefooter  No aparece
	 * TableHeader
	 * TableRow
*/
type Tabla struct {
	NodoHeader
}

func NewTabla() Tabla {
	return Tabla{
		NodoHeader: NewTipo(MK_Table),
	}
}
