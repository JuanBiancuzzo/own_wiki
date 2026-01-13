package database

type FieldType uint8

const (
	FT_INT = iota
	FT_STRING
	FT_CHAR
	FT_BOOL
	FT_DATE
	FT_REF
)

type TableDescription struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name   string
	Type   FieldType
	IsNull bool
	IsKey  bool
}
