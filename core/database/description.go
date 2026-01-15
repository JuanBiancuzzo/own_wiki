package database

// This Could be a primitive value, int, string, etc. or a pointer to
// another TablaData
type FieldData any

type DataAmount uint8

const (
	DA_POINT = iota
	DA_ARRAY
)

type TableData struct {
	Structure TableStructure

	DataAmount

	// This is a matrix of data, where the data is organize a continuos concatenation
	// of rows, that are the data of the table. The length of the rows is given by the
	// length of the fields in the TableStructure
	Data []FieldData
}

// This TableElement could be a TableData o TableComposition
type TableElement any

type TableComposition struct {
	Composition []TableElement
}
