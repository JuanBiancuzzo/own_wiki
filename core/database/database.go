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

func NewTableDescription(name string, fields []Field) *TableDescription {
	return &TableDescription{
		Name:   name,
		Fields: fields,
	}
}

type Field struct {
	Name      string
	Type      FieldType
	Reference *TableDescription
	IsNull    bool
	IsKey     bool
}

func NewPrimitiveField(name string, fieldType FieldType, isNull, isKey bool) Field {
	return Field{
		Name:      name,
		Type:      fieldType,
		Reference: nil,
		IsNull:    isNull,
		IsKey:     isKey,
	}
}

func NewReferencesField(name string, reference *TableDescription, isNull, isKey bool) Field {
	return Field{
		Name:      name,
		Type:      FT_REF,
		Reference: reference,
		IsNull:    isNull,
		IsKey:     isKey,
	}
}
