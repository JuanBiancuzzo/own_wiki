package estructura

type Dependencia struct {
}

/*
	una dependencia apunta a otra dependencia siguiente en la cadena
	una dependencia apunta a una funcion que se corre con el resultado de la cadena de dependencias

	dependencia {
		tabla
		info query (aka valor buscado)
		valorGuardado []any
		dependenciaCadena []*dependencia
		cargable *Cargable
	}
*/
