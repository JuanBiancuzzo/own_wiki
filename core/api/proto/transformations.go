package core

import (
	"fmt"

	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

func (ft PrimitiveFieldType) GetDataBaseFieldType() (fieldType db.FieldType, err error) {
	switch ft {
	case PrimitiveFieldType_INT:
		fieldType = db.FT_INT

	case PrimitiveFieldType_STRING:
		fieldType = db.FT_STRING

	case PrimitiveFieldType_CHAR:
		fieldType = db.FT_CHAR

	case PrimitiveFieldType_BOOL:
		fieldType = db.FT_BOOL

	case PrimitiveFieldType_DATE:
		fieldType = db.FT_DATE

	default:
		err = fmt.Errorf("Field %d not define", ft)
	}

	return fieldType, err
}

func GetFromDataBaseFieldType(ft db.FieldType) (fieldType PrimitiveFieldType, err error) {
	switch ft {
	case db.FT_INT:
		fieldType = PrimitiveFieldType_INT

	case db.FT_STRING:
		fieldType = PrimitiveFieldType_STRING

	case db.FT_CHAR:
		fieldType = PrimitiveFieldType_CHAR

	case db.FT_BOOL:
		fieldType = PrimitiveFieldType_BOOL

	case db.FT_DATE:
		fieldType = PrimitiveFieldType_DATE

	default:
		err = fmt.Errorf("Field %d is not define", ft)
	}

	return fieldType, err
}

func GetFromDataBaseDataAmount(da db.DataAmount) (amount ComponentDescription_DataAmount, err error) {
	switch da {
	case db.DA_POINT:
		amount = ComponentDescription_ONE

	case db.DA_ARRAY:
		amount = ComponentDescription_ARRAY

	default:
		err = fmt.Errorf("Data amount %d is not define", da)
	}

	return amount, err
}

func GetFieldData(fieldType db.FieldType, data any) (dataValue *FieldDescription_FieldData, err error) {
	dataValue = &FieldDescription_FieldData{}

	switch fieldType {
	case db.FT_INT:
		if number, ok := data.(int); !ok {
			err = fmt.Errorf("data value is not an int, and as a number it should be")

		} else {
			dataValue.Data = &FieldDescription_FieldData_Number{Number: int32(number)}
		}

	case db.FT_STRING:
		if text, ok := data.(string); !ok {
			err = fmt.Errorf("data value is not an string, and as a text it should be")

		} else {
			dataValue.Data = &FieldDescription_FieldData_Text{Text: text}
		}

	case db.FT_CHAR:
		if character, ok := data.(string); !ok {
			err = fmt.Errorf("data value is not an string, and as a character it should be")

		} else {
			dataValue.Data = &FieldDescription_FieldData_Character{Character: character}
		}

	case db.FT_BOOL:
		if boolean, ok := data.(bool); !ok {
			err = fmt.Errorf("data value is not an bool, and as a boolean it should be")

		} else {
			dataValue.Data = &FieldDescription_FieldData_Boolean{Boolean: boolean}
		}

	case db.FT_DATE:
		if date, ok := data.(uint); !ok {
			err = fmt.Errorf("data value is not an uint, and as a date it should be")

		} else {
			dataValue.Data = &FieldDescription_FieldData_Date{Date: uint32(date)}
		}

	case db.FT_REF:
		if reference, ok := data.(*db.TableData); !ok {
			err = fmt.Errorf("data value is not an pointer to the table data, and as a reference it should be")

		} else if refComponent, errComp := ConvertToComponentDescription(reference); errComp != nil {
			err = fmt.Errorf("failed to create reference component, with error: %v", errComp)

		} else {
			dataValue.Data = &FieldDescription_FieldData_Reference{
				Reference: refComponent.Component,
			}
		}

	default:
		err = fmt.Errorf("field %d is not define", fieldType)
	}

	return dataValue, err
}

func ConvertToComponentDescription(tableData *db.TableData) (*EntityDescription_Component, error) {
	amountFields := len(tableData.Structure.Fields)
	fields := make([]*FieldDescription, amountFields)

	dataAmount, err := GetFromDataBaseDataAmount(tableData.DataAmount)
	if err != nil {
		return nil, fmt.Errorf("invalid tablaData amount, with error: %v", err)
	}

	dataRows := len(tableData.Data) / amountFields

	for i, fieldData := range tableData.Structure.Fields {
		var fieldTypeInformation *FieldTypeInformation
		switch fieldData.Type {
		case db.FT_INT, db.FT_STRING, db.FT_CHAR, db.FT_BOOL, db.FT_DATE:
			primitiveType, err := GetFromDataBaseFieldType(fieldData.Type)
			if err != nil {
				return nil, fmt.Errorf("invalid type of primitive in %s field, with error: %v", fieldData.Name, err)
			}

			fieldTypeInformation = &FieldTypeInformation{
				Type: FieldTypeInformation_PRIMITIVE,
				Information: &FieldTypeInformation_Primitive{
					Primitive: primitiveType,
				},
			}

		case db.FT_REF:
			fieldTypeInformation = &FieldTypeInformation{
				Type: FieldTypeInformation_REFERENCE,
				Information: &FieldTypeInformation_Reference{
					Reference: &ReferenceInformation{
						TableName: fieldData.Reference.Name,
					},
				},
			}

		default:
			return nil, fmt.Errorf("field %s has invalid type of %v", fieldData.Name, fieldData.Type)
		}

		fields[i] = &FieldDescription{
			Name:            fieldData.Name,
			TypeInformation: fieldTypeInformation,
		}

		switch tableData.DataAmount {
		case db.DA_POINT:
			dataValue, err := GetFieldData(fieldData.Type, tableData.Data[i])
			if err != nil {
				return nil, fmt.Errorf("field data in invalid, with error: %v", err)
			}

			fields[i].Data = &FieldDescription_Point{Point: dataValue}

		case db.DA_ARRAY:
			dataArray := make([]*FieldDescription_FieldData, dataRows)
			for row := range dataRows {
				value := tableData.Data[i+row*amountFields]
				if dataValue, err := GetFieldData(fieldData.Type, value); err != nil {
					return nil, fmt.Errorf("field data in invalid, with error: %v", err)

				} else {
					dataArray[row] = dataValue
				}
			}

			fields[i].Data = &FieldDescription_Array{
				Array: &FieldDescription_FieldDataArray{DataArray: dataArray},
			}
		}
	}

	return &EntityDescription_Component{
		Component: &ComponentDescription{
			Name:   tableData.Structure.Name,
			Amount: dataAmount,
			Fields: fields,
		},
	}, nil
}

func ConvertToCompositionDescription(tableComposition *db.TableComposition) (*EntityDescription_Composition, error) {
	entities := make([]*EntityDescription, len(tableComposition.Composition))

	for i, tableElement := range tableComposition.Composition {
		if entity, err := ConvertToEntityDescription(&tableElement); err != nil {
			return nil, fmt.Errorf("Failed to convert entity description while converting from table composition, with error: %v", err)

		} else {
			entities[i] = entity
		}
	}

	return &EntityDescription_Composition{
		Composition: &ComponentCompositionDescription{
			Entities: entities,
		},
	}, nil
}

func ConvertToEntityDescription(tableElement *db.TableElement) (*EntityDescription, error) {
	switch table := (*tableElement).(type) {
	case db.TableData:
		if component, err := ConvertToComponentDescription(&table); err != nil {
			return nil, err

		} else {
			return &EntityDescription{Description: component}, nil
		}

	case db.TableComposition:
		if composition, err := ConvertToCompositionDescription(&table); err != nil {
			return nil, err

		} else {
			return &EntityDescription{Description: composition}, nil
		}

	default:
		return nil, fmt.Errorf("Table element is not a data o composition type")
	}
}
