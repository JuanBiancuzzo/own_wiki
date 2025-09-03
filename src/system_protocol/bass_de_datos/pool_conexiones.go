package bass_de_datos

import "fmt"

const MAX_CONEXIONES = 50
const CONEXIONES_INICIALES = 5

type poolConexiones struct {
	archivo            string
	cantidadConexiones int
	conexiones         chan *conexion
}

func newPoolConexiones(archivoBdd string) (*poolConexiones, error) {
	conexiones := make(chan *conexion, MAX_CONEXIONES)
	for range CONEXIONES_INICIALES {
		if conn, err := newConexion(archivoBdd, conexiones); err != nil {
			return nil, err
		} else {
			conexiones <- conn
		}
	}

	return &poolConexiones{
		archivo:            archivoBdd,
		cantidadConexiones: CONEXIONES_INICIALES,
		conexiones:         conexiones,
	}, nil
}

func (pc *poolConexiones) Conexion() (*conexion, error) {
	select {
	case conn := <-pc.conexiones:
		return conn, nil

	default:
		if pc.cantidadConexiones >= MAX_CONEXIONES {
			fmt.Println("A Esperar que una conexion se libere")
			return <-pc.conexiones, nil
		}

		if conn, err := newConexion(pc.archivo, pc.conexiones); err != nil {
			return nil, err
		} else {
			fmt.Println("Creando nueva conexionnnn")
			pc.cantidadConexiones++
			return conn, nil
		}
	}
}

func (pc *poolConexiones) Close() {
	for range pc.cantidadConexiones {
		conn := <-pc.conexiones
		conn.Close()
	}
}
