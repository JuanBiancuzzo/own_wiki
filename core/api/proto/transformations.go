package core

import (
	"fmt"
	"time"

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
	switch fieldType {
	case db.FT_INT:
		if number, ok := data.(int); ok {
			dataValue = NewFieldConcreteNumber(number)

		} else if number, ok := data.(*int); ok {
			dataValue = NewFieldNullableNumber(number)

		} else {
			err = fmt.Errorf("data value is not an int, and as a number it should be")
		}

	case db.FT_STRING:
		if text, ok := data.(string); ok {
			dataValue = NewFieldConcreteText(text)

		} else if text, ok := data.(*string); ok {
			dataValue = NewFieldNullableText(text)

		} else {
			err = fmt.Errorf("data value is not an string, and as a text it should be")
		}

	case db.FT_CHAR:
		if character, ok := data.(string); ok {
			dataValue = NewFieldConcreteCharacter(character)

		} else if character, ok := data.(*string); ok {
			dataValue = NewFieldNullableCharacter(character)

		} else {
			err = fmt.Errorf("data value is not an string, and as a character it should be")
		}

	case db.FT_BOOL:
		if boolean, ok := data.(bool); ok {
			dataValue = NewFieldConcreteBoolean(boolean)

		} else if boolean, ok := data.(*bool); ok {
			dataValue = NewFieldNullableBoolean(boolean)

		} else {
			err = fmt.Errorf("data value is not an bool, and as a boolean it should be")
		}

	case db.FT_DATE:
		if date, ok := data.(time.Time); ok {
			dataValue = NewFieldConcreteDate(date)

		} else if date, ok := data.(*time.Time); ok {
			dataValue = NewFieldNullableDate(date)

		} else {
			err = fmt.Errorf("data value is not an time.Time, and as a date it should be")
		}

	case db.FT_REF:
		// TODO: make reference

	default:
		err = fmt.Errorf("field %d is not define", fieldType)
	}

	return dataValue, err
}

func ConvertToComponentDescription(tableData *db.TableData) (*ComponentDescription, error) {
	amountFields := len(tableData.Structure.Fields)
	fields := make([]*FieldDescription, amountFields)

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
			fieldTypeInformation = NewCompletePrimitiveInfo(primitiveType)

		case db.FT_REF:
			fieldTypeInformation = NewReferenceInfo(fieldData.Reference.Name, fieldData.IsNull)

		default:
			return nil, fmt.Errorf("field %s has invalid type of %v", fieldData.Name, fieldData.Type)
		}

		switch tableData.DataAmount {
		case db.DA_POINT:
			dataValue, err := GetFieldData(fieldData.Type, tableData.Data[i])
			if err != nil {
				return nil, fmt.Errorf("field data in invalid, with error: %v", err)
			}

			fields[i] = NewFieldPoint(fieldData.Name, fieldTypeInformation, dataValue)

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

			fields[i] = NewFieldArray(fieldData.Name, fieldTypeInformation, dataArray)

		default:
			return nil, fmt.Errorf("invalid tablaData amount")
		}
	}

	switch tableData.DataAmount {
	case db.DA_POINT:
		return NewComponentDescriptionPoint(tableData.Structure.Name, fields), nil

	case db.DA_ARRAY:
		return NewComponentDescriptionArray(tableData.Structure.Name, int32(dataRows), fields), nil

	default:
		panic("This should be possible")
	}
}

func ConvertToCompositionDescription(tableComposition *db.TableComposition) (*ComponentCompositionDescription, error) {
	entities := make([]*EntityDescription, len(tableComposition.Composition))

	for i, tableElement := range tableComposition.Composition {
		if entity, err := ConvertToEntityDescription(&tableElement); err != nil {
			return nil, fmt.Errorf("Failed to convert entity description while converting from table composition, with error: %v", err)

		} else {
			entities[i] = entity
		}
	}

	return &ComponentCompositionDescription{
		Entities: entities,
	}, nil
}

func ConvertToEntityDescription(tableElement *db.TableElement) (entity *EntityDescription, errEntity error) {
	switch table := (*tableElement).(type) {
	case *db.TableData:
		if component, err := ConvertToComponentDescription(table); err != nil {
			errEntity = err

		} else {
			entity = NewEntityDescriptionComponent(component)
		}

	case *db.TableComposition:
		if composition, err := ConvertToCompositionDescription(table); err != nil {
			errEntity = err

		} else {
			entity = NewEntityDescriptionComposition(composition)
		}

	default:
		errEntity = fmt.Errorf("Table element is not a data o composition type")
	}

	return entity, errEntity
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
			data = getTypeValue(isNull, dataArray, func(concrete *FieldDescription_ConcreteFieldData) time.Time {
				return time.Unix(int64(concrete.GetDate()), 0)
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
