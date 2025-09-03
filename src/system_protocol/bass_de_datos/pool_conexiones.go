package bass_de_datos

import "sync"

const MAX_CONEXIONES = 50
const CONEXIONES_INICIALES = 5

type poolConexiones struct {
	archivo            string
	cantidadConexiones int
	conexiones         chan *conexion

	waitGroup *sync.WaitGroup
	esperando bool
}

func newPoolConexiones(archivoBdd string) (*poolConexiones, error) {
	conexiones := make(chan *conexion, MAX_CONEXIONES)
	for _ = range CONEXIONES_INICIALES {
		if conn, err := newConexion(archivoBdd, conexiones); err != nil {
			return nil, err
		} else {
			conexiones <- conn
		}
	}

	var waitGroup sync.WaitGroup
	return &poolConexiones{
		archivo:            archivoBdd,
		cantidadConexiones: CONEXIONES_INICIALES,
		conexiones:         conexiones,

		waitGroup: &waitGroup,
		esperando: false,
	}, nil
}

func (pc *poolConexiones) Conexion() (*conexion, error) {
	select {
	case conn := <-pc.conexiones:
		return conn, nil

	default:
		if pc.cantidadConexiones >= MAX_CONEXIONES {
			return <-pc.conexiones, nil
		}

		if conn, err := newConexion(pc.archivo, pc.conexiones); err != nil {
			return nil, err
		} else {
			pc.cantidadConexiones++
			return conn, nil
		}
	}
}

func (pc *poolConexiones) Close() {
	for _ = range pc.cantidadConexiones {
		conn := <-pc.conexiones
		conn.Close()
	}
}
