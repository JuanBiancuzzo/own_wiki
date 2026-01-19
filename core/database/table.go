package database

type FieldType uint8

const (
	FT_INT = iota
	FT_STRING
	FT_BOOL
	FT_DATE
	FT_REF
)

type TableStructure struct {
	Name   string
	Fields []Field
}

func NewTableStructure(name string, fields []Field) *TableStructure {
	return &TableStructure{
		Name:   name,
		Fields: fields,
	}
}

func (ts *TableStructure) AmountOfPrimitiveValues() (amount uint) {
	for _, field := range ts.Fields {
		amount += field.AmountOfPrimitiveValues()
	}
	return amount
}

type Field struct {
	Name      string
	Type      FieldType
	Reference *TableStructure
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

func NewReferencesField(name string, reference *TableStructure, isNull, isKey bool) Field {
	return Field{
		Name:      name,
		Type:      FT_REF,
		Reference: reference,
		IsNull:    isNull,
		IsKey:     isKey,
	}
}

func (f Field) AmountOfPrimitiveValues() (amount uint) {
	switch f.Type {
	case FT_INT, FT_STRING, FT_BOOL, FT_DATE:
		amount = 1

	case FT_REF:
		amount = 1
		// amount = f.Reference.AmountOfPrimitiveValues()
	}
	return amount
}
