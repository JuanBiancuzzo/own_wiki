package database

type FieldData any

type TableData struct {
	Structure TableStructure

	// This is a matrix of data, where the data is organize a continuos concatenation
	// of rows, that are the data of the table. The length of the rows is given by the
	// length of the fields in the TableStructure
	Data []FieldData
}
