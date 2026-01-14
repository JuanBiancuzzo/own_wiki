package core

import (
	"fmt"

	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

func (ft PrimiteFieldType) GetDataBaseFieldType() (fieldType db.FieldType, err error) {
	switch ft {
	case PrimiteFieldType_INT:
		fieldType = db.FT_INT

	case PrimiteFieldType_STRING:
		fieldType = db.FT_STRING

	case PrimiteFieldType_CHAR:
		fieldType = db.FT_CHAR

	case PrimiteFieldType_BOOL:
		fieldType = db.FT_BOOL

	case PrimiteFieldType_DATE:
		fieldType = db.FT_DATE

	default:
		return fieldType, fmt.Errorf("Field %d not define", ft)
	}

	return fieldType, nil
}
