package database

// This could be a primitive value, int, string, etc. Note: only primitive value and
// there pointers could be FieldData. There shouldn't be a pointer to a structure
type FieldData any

// This TableElement could be a *TableData o *TableComposition
type TableElement any

type DataAmount uint8

const (
	DA_POINT = iota
	DA_ARRAY
)

type TableData struct {
	Structure *TableStructure

	DataAmount

	// This is a matrix of data, where the data is organize a continuos concatenation
	// of rows, that are the data of the table. The length of the rows is given by the
	// length of the fields in the TableStructure
	Data []FieldData
}

func NewTableData(structure *TableStructure, dataAmount DataAmount, data []FieldData) *TableData {
	return &TableData{
		Structure:  structure,
		DataAmount: dataAmount,
		Data:       data,
	}
}

type TableComposition struct {
	Composition []TableElement
}

func NewTableComposition(composition ...TableElement) *TableComposition {
	return &TableComposition{
		Composition: composition,
	}
}
