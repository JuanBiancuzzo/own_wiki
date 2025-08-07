package utilidades

type InfiniteWriter struct {
	bytes *Lista[byte]
}

func NewInfiniteWriter() *InfiniteWriter {
	return &InfiniteWriter{
		bytes: NewLista[byte](),
	}
}

func (iw *InfiniteWriter) Write(p []byte) (n int, err error) {
	n = 0
	for _, byteAgregar := range p {
		iw.bytes.Push(byteAgregar)
		n++
	}

	return n, nil
}

func (iw *InfiniteWriter) Reset() {
	iw.bytes.Vaciar()
}

func (iw *InfiniteWriter) Items() []byte {
	return iw.bytes.Items()
}
