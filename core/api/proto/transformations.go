package core

import (
	"fmt"

	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

func (ft FieldType) GetDataBaseFieldType() (fieldType db.FieldType, err error) {
	switch ft {
	case FieldType_INT:
		fieldType = db.FT_INT
	case FieldType_STRING:
		fieldType = db.FT_STRING
	case FieldType_CHAR:
		fieldType = db.FT_CHAR
	case FieldType_BOOL:
		fieldType = db.FT_BOOL
	case FieldType_DATE:
		fieldType = db.FT_DATE

	default:
		return fieldType, fmt.Errorf("Field %d not define", ft)
	}

	return fieldType, nil
}
