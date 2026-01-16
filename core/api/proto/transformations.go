package core

import (
	"fmt"

	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

func (ft PrimitiveFieldType) GetDataBaseFieldType() (fieldType db.FieldType, isNull bool, err error) {
	isNull = ft >= PrimitiveFieldType_NULL_INT
	if isNull {
		ft -= PrimitiveFieldType_NULL_INT
	}

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

	return fieldType, isNull, err
}

func GetFromDataBaseFieldType(ft db.FieldType, isNull bool) (fieldType PrimitiveFieldType, err error) {
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

	if isNull {
		// As the NULL_INT is the first nullable element, then all the next ones are ofset by the db.FieldType
		fieldType += PrimitiveFieldType_NULL_INT
	}

	return fieldType, err
}

func (da ComponentDescription_DataAmount) GetDataBaseAmount() (dataAmount db.DataAmount, err error) {
	switch da {
	case ComponentDescription_ONE:
		dataAmount = db.DA_POINT

	case ComponentDescription_ARRAY:
		dataAmount = db.DA_ARRAY

	default:
		err = fmt.Errorf("Data amount %d is not define", da)
	}

	return dataAmount, err
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
	value := &FieldDescription_ConcreteFieldData{}

	isConcrete := false

	switch fieldType {
	case db.FT_INT:
		if number, ok := data.(int); ok {
			value.Data = &FieldDescription_ConcreteFieldData_Number{Number: int32(number)}
			isConcrete = true

		} else if number, ok := data.(*int); !ok {
			err = fmt.Errorf("data value is not an int, and as a number it should be")

		} else if number != nil {
			value.Data = &FieldDescription_ConcreteFieldData_Number{Number: int32(*number)}
		}

	case db.FT_STRING:
		if text, ok := data.(string); ok {
			value.Data = &FieldDescription_ConcreteFieldData_Text{Text: text}
			isConcrete = true

		} else if text, ok := data.(*string); !ok {
			err = fmt.Errorf("data value is not an string, and as a text it should be")

		} else if text != nil { // its a *string and its not nil, then the value.Data is not nil
			value.Data = &FieldDescription_ConcreteFieldData_Text{Text: *text}
		}

	case db.FT_CHAR:
		if character, ok := data.(string); ok {
			value.Data = &FieldDescription_ConcreteFieldData_Character{Character: character}
			isConcrete = true

		} else if character, ok := data.(*string); !ok {
			err = fmt.Errorf("data value is not an string, and as a character it should be")

		} else if character != nil {
			value.Data = &FieldDescription_ConcreteFieldData_Character{Character: *character}
		}

	case db.FT_BOOL:
		if boolean, ok := data.(bool); ok {
			value.Data = &FieldDescription_ConcreteFieldData_Boolean{Boolean: boolean}
			isConcrete = true

		} else if boolean, ok := data.(*bool); !ok {
			err = fmt.Errorf("data value is not an bool, and as a boolean it should be")

		} else if boolean != nil {
			value.Data = &FieldDescription_ConcreteFieldData_Boolean{Boolean: *boolean}
		}

	case db.FT_DATE:
		if date, ok := data.(uint); ok {
			value.Data = &FieldDescription_ConcreteFieldData_Date{Date: uint32(date)}
			isConcrete = true

		} else if date, ok := data.(*uint); !ok {
			err = fmt.Errorf("data value is not an uint, and as a date it should be")

		} else if date != nil {
			value.Data = &FieldDescription_ConcreteFieldData_Date{Date: uint32(*date)}
		}

	case db.FT_REF:
		if reference, ok := data.(*db.TableData); !ok {
			err = fmt.Errorf("data value is not an pointer to the table data, and as a reference it should be")

		} else if reference == nil {

		} else if refComponent, errComp := ConvertToComponentDescription(reference); errComp != nil {
			err = fmt.Errorf("failed to create reference component, with error: %v", errComp)

		} else {
			value.Data = &FieldDescription_ConcreteFieldData_Reference{
				Reference: refComponent.Component,
			}
		}

	default:
		err = fmt.Errorf("field %d is not define", fieldType)
	}

	if isConcrete {
		dataValue = &FieldDescription_FieldData{
			Data: &FieldDescription_FieldData_Concrete{Concrete: value},
		}

	} else if value.Data != nil {
		dataValue = &FieldDescription_FieldData{
			Data: &FieldDescription_FieldData_Nullable{
				Nullable: &FieldDescription_NullableFieldData{Data: value},
			},
		}

	} else {
		dataValue = &FieldDescription_FieldData{
			Data: &FieldDescription_FieldData_Nullable{
				Nullable: &FieldDescription_NullableFieldData{Data: nil},
			},
		}
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

	dataColumns := tableData.Structure.AmountOfPrimitiveValues()
	dataRows := len(tableData.Data) / int(dataColumns)

	for i, fieldData := range tableData.Structure.Fields {
		var fieldTypeInformation *FieldTypeInformation
		switch fieldData.Type {
		case db.FT_INT, db.FT_STRING, db.FT_CHAR, db.FT_BOOL, db.FT_DATE:
			primitiveType, err := GetFromDataBaseFieldType(fieldData.Type, fieldData.IsNull)
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
					return nil, fmt.Errorf("field data is invalid, with error: %v", err)

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
			Rows:   int32(dataRows),
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
	case *db.TableData:
		if component, err := ConvertToComponentDescription(table); err != nil {
			return nil, err

		} else {
			return &EntityDescription{Description: component}, nil
		}

	case *db.TableComposition:
		if composition, err := ConvertToCompositionDescription(table); err != nil {
			return nil, err

		} else {
			return &EntityDescription{Description: composition}, nil
		}

	default:
		return nil, fmt.Errorf("Table element is not a data o composition type")
	}
}

type fnGetValue[T any] func(*FieldDescription_ConcreteFieldData) T

func getTypeValue[T any](isNull bool, dataArray []*FieldDescription_FieldData, getValue fnGetValue[T]) []db.FieldData {
	data := make([]db.FieldData, len(dataArray))

	if isNull {
		for i, dataValue := range dataArray {
			var x *T = nil
			if concrete := dataValue.GetNullable().GetData(); concrete != nil {
				value := getValue(concrete)
				x = &value
			}
			data[i] = x
		}

	} else {
		for i, dataValue := range dataArray {
			data[i] = getValue(dataValue.GetConcrete())
		}

		return data
	}

	return data
}

func (fd *FieldDescription) ConvertToDataValues(dataAmount db.DataAmount) (data []db.FieldData, field db.Field, err error) {
	typeInformation := fd.GetTypeInformation()
	var dataArray []*FieldDescription_FieldData
	switch dataAmount {
	case db.DA_POINT:
		dataArray = []*FieldDescription_FieldData{fd.GetPoint()}

	case db.DA_ARRAY:
		dataArray = fd.GetArray().GetDataArray()
	}

	switch typeInformation.GetType() {
	case FieldTypeInformation_PRIMITIVE:
		fieldType, isNull, errFT := typeInformation.GetPrimitive().GetDataBaseFieldType()
		if errFT != nil {
			err = fmt.Errorf("Failed to get field type, with error: %v", err)
			break
		}

		switch fieldType {
		case db.FT_INT:
			data = getTypeValue(isNull, dataArray, func(concrete *FieldDescription_ConcreteFieldData) int {
				return int(concrete.GetNumber())
			})

		case db.FT_STRING:
			data = getTypeValue(isNull, dataArray, func(concrete *FieldDescription_ConcreteFieldData) string {
				return concrete.GetText()
			})

		case db.FT_CHAR:
			data = getTypeValue(isNull, dataArray, func(concrete *FieldDescription_ConcreteFieldData) string {
				return concrete.GetCharacter()
			})

		case db.FT_BOOL:
			data = getTypeValue(isNull, dataArray, func(concrete *FieldDescription_ConcreteFieldData) bool {
				return concrete.GetBoolean()
			})

		case db.FT_DATE:
			data = getTypeValue(isNull, dataArray, func(concrete *FieldDescription_ConcreteFieldData) uint {
				return uint(concrete.GetDate())
			})
		}

		// The literal false is for the isKey, that for the date is irrelevant
		field = db.NewPrimitiveField(fd.GetName(), fieldType, isNull, false)

	/*
		case FieldTypeInformation_REFERENCE:
			switch value := dataPoint.GetData().(type) {
			case *FieldDescription_FieldData_Nullable:
				if concrete := value.Nullable.GetData(); concrete == nil {
					structure := db.TableStructure{
						Name: typeInformation.GetReference().GetTableName(),
					}

					var table *db.TableData = nil
					data = table
					field = db.NewReferencesField(fd.GetName(), &structure, true, false)

				} else if tableData, errRef := concrete.GetReference().ConvertToTableData(); errRef == nil {
					data = &tableData
					field = db.NewReferencesField(fd.GetName(), &tableData.Structure, false, false)

				} else {
					// TODO: handle error
					err = errRef
				}

			case *FieldDescription_FieldData_Concrete:
				if tableData, errRef := value.Concrete.GetReference().ConvertToTableData(); errRef == nil {
					data = tableData
					field = db.NewReferencesField(fd.GetName(), &tableData.Structure, false, false)

				} else {
					// TODO: handle error
					err = errRef
				}
			}
	*/

	default:
		err = fmt.Errorf("Invalid type information, with %d", typeInformation.GetType())
	}

	return data, field, err
}

func (cd *ComponentDescription) ConvertToTableData() (*db.TableData, error) {
	dataAmount, err := cd.GetAmount().GetDataBaseAmount()
	if err != nil {
		return nil, fmt.Errorf("Failed to get dataAmount, with error: %v", err)
	}

	amountFields := len(cd.GetFields())
	data := make([]db.FieldData, int(cd.GetRows())*amountFields)
	fields := make([]db.Field, amountFields)

	// As each field may contain more than one element, this should acount for the offset of the data of each field
	acumulatedOffset := 0
	for i, fieldDescription := range cd.GetFields() {
		dataArray, field, err := fieldDescription.ConvertToDataValues(dataAmount)
		if err != nil {
			return nil, fmt.Errorf("Failed to get DataArray data and field infomation, with error: %v", err)

		} else if len(dataArray) != int(cd.GetRows()) {
			return nil, fmt.Errorf("The amount of data in the arrays (%d) is not the same as the amount of rows (%d)", len(dataArray), cd.Rows)
		}
		fields[i] = field

		// The offset should be the amount of elements data are store
		addedOffet := int(field.AmountOfPrimitiveValues())

		for row := range len(dataArray) {
			for offset := range addedOffet {
				dataPoint := dataArray[row+offset]
				data[offset+acumulatedOffset+row*amountFields] = dataPoint
			}
		}
		acumulatedOffset += addedOffet
	}

	structure := db.NewTableStructure(cd.GetName(), fields)
	return db.NewTableData(structure, dataAmount, data), nil
}

func (ed *EntityDescription) ConvertToTableElement() (tableElement db.TableElement, err error) {
	switch table := ed.GetDescription().(type) {
	case *EntityDescription_Component:
		tableElement, err = table.Component.ConvertToTableData()

	case *EntityDescription_Composition:
		// TODO: Convert composition

	default:
		err = fmt.Errorf("The entity description is not define")
	}

	return tableElement, err
}
