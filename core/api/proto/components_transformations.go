package core

import (
	"fmt"
	"time"

	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

type CacheStructures struct {
	TableStructures map[string]*db.TableStructure
}

func NewCacheData() *CacheStructures {
	return &CacheStructures{
		TableStructures: make(map[string]*db.TableStructure),
	}
}

func (ft FieldType) ConvertToDatabase() (fieldType db.FieldType, err error) {
	switch ft {
	case FieldType_INT:
		fieldType = db.FT_INT

	case FieldType_STRING:
		fieldType = db.FT_STRING

	case FieldType_BOOL:
		fieldType = db.FT_BOOL

	case FieldType_DATE:
		fieldType = db.FT_DATE

	case FieldType_REFERENCE:
		fieldType = db.FT_REF

	default:
		err = fmt.Errorf("Field %d not define", ft)
	}

	return fieldType, err
}

func GetFromDataBaseFieldType(ft db.FieldType) (fieldType FieldType, err error) {
	switch ft {
	case db.FT_INT:
		fieldType = FieldType_INT

	case db.FT_STRING:
		fieldType = FieldType_STRING

	case db.FT_BOOL:
		fieldType = FieldType_BOOL

	case db.FT_DATE:
		fieldType = FieldType_DATE

	default:
		err = fmt.Errorf("Field %d is not define", ft)
	}

	return fieldType, err
}

func (da DataAmount) GetDataBaseAmount() (dataAmount db.DataAmount, err error) {
	switch da {
	case DataAmount_POINT:
		dataAmount = db.DA_POINT

	case DataAmount_ARRAY:
		dataAmount = db.DA_ARRAY

	default:
		err = fmt.Errorf("Data amount %d is not define", da)
	}

	return dataAmount, err
}

func GetFromDataBaseDataAmount(da db.DataAmount) (amount DataAmount, err error) {
	switch da {
	case db.DA_POINT:
		amount = DataAmount_POINT

	case db.DA_ARRAY:
		amount = DataAmount_ARRAY

	default:
		err = fmt.Errorf("Data amount %d is not define", da)
	}

	return amount, err
}

func ConvertToFieldStructure(field *db.Field) (fieldStructure *FieldStructure, err error) {
	switch field.Type {
	case db.FT_INT, db.FT_STRING, db.FT_BOOL, db.FT_DATE:
		var fieldType FieldType
		if fieldType, err = GetFromDataBaseFieldType(field.Type); err == nil {
			fieldStructure = NewPrimitiveStructure(field.Name, fieldType, field.IsNull, field.IsKey)
		}

	case db.FT_REF:
		var componentStructure *ComponentStructure
		if componentStructure, err = ConvertToComponentStructure(field.Reference); err == nil {
			fieldStructure = NewReferenceStructure(field.Name, componentStructure, field.IsNull, field.IsKey)
		}

	default:
		err = fmt.Errorf("Invalid field type of: %d", field.Type)
	}

	return fieldStructure, err
}

func ConvertToComponentStructure(structure *db.TableStructure) (componentStructure *ComponentStructure, err error) {
	fields := make([]*FieldStructure, len(structure.Fields))
	for i, field := range structure.Fields {
		if fields[i], err = ConvertToFieldStructure(&field); err != nil {
			return componentStructure, err
		}
	}
	return NewComponentStructure(structure.Name, fields), err
}

func ConvertFieldDescriptions(structure *db.TableStructure, data []db.FieldData) (fieldDescriptions []*FieldDescription, rows int, err error) {
	columns := len(structure.Fields)
	rows = len(data) / columns

	fieldDescriptions = make([]*FieldDescription, columns)
	for column, field := range structure.Fields {
		switch field.Type {
		case db.FT_INT, db.FT_STRING, db.FT_BOOL, db.FT_DATE:
			verticalData := make([]any, rows)
			for row := range rows {
				verticalData[row] = data[column+row*columns]
			}

			var fieldData []*FieldData
			if fieldData, err = NewFieldDataArrayWithError(verticalData); err == nil {
				fieldDescriptions[column] = NewPrimitiveDescription(fieldData...)
			}

		case db.FT_REF:
			refStructure := field.Reference
			refColumns := len(refStructure.Fields)
			// We should get the data in an array from, concatenated instead of an array of array
			refData := make([]db.FieldData, rows*refColumns)
			for row := range rows {
				for refColumn := range refColumns {
					if dataArray, ok := data[column+row*columns].([]db.FieldData); ok {
						refData[refColumn+row*refColumns] = dataArray[refColumn]
					} else {
						return fieldDescriptions, rows, fmt.Errorf("References date not found, the data was of type %T", data[column+row*columns])
					}
				}
			}

			var refDescriptions []*FieldDescription
			if refDescriptions, _, err = ConvertFieldDescriptions(refStructure, refData); err == nil {
				fieldDescriptions[column] = NewReferenceDescription(refDescriptions...)
			}

		default:
			err = fmt.Errorf("No type of value: %d", field.Type)
		}
	}

	return fieldDescriptions, rows, err
}

func ConvertToComponentDescription(tableData *db.TableData) (component *ComponentDescription, errComp error) {
	if structure, err := ConvertToComponentStructure(tableData.Structure); err != nil {
		errComp = fmt.Errorf("Failed to convert structure data, with error: %v", err)

	} else if fieldDescriptions, rows, err := ConvertFieldDescriptions(tableData.Structure, tableData.Data); err != nil {
		errComp = fmt.Errorf("Failed to convert data in fieldDescriptions, with error: %v", err)

	} else if tableData.DataAmount == db.DA_POINT {
		component = NewComponentDescriptionPoint(structure, fieldDescriptions...)

	} else if tableData.DataAmount == db.DA_ARRAY {
		component = NewComponentDescriptionArray(structure, rows, fieldDescriptions...)

	} else {
		errComp = fmt.Errorf("Invalid data amount of %d", tableData.DataAmount)
	}

	return component, errComp
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

func getTypeValue[T any](dataArray []*FieldData, isNullable bool) []db.FieldData {
	data := make([]db.FieldData, len(dataArray))
	for i, dataPoint := range dataArray {
		var getValue func(*FieldData) db.FieldData
		var x T
		switch any(x).(type) {
		case int, *int:
			getValue = func(fd *FieldData) db.FieldData { return int(fd.GetNumber()) }

		case string, *string:
			getValue = func(fd *FieldData) db.FieldData { return fd.GetText() }

		case bool, *bool:
			getValue = func(fd *FieldData) db.FieldData { return fd.GetBoolean() }

		case time.Time, *time.Time:
			getValue = func(fd *FieldData) db.FieldData {
				return time.Unix(int64(fd.GetDate()), 0)
			}
		}

		if isNullable && dataPoint.IsNull {
			var nilValue *T = nil
			data[i] = nilValue

		} else {
			concreteValue := getValue(dataPoint)
			if value, ok := concreteValue.(T); ok && isNullable {
				data[i] = &value
			} else if ok {
				data[i] = value
			}
		}
	}

	return data
}

func (fd *FieldDescription) ConvertToDataValues(fieldStructure *FieldStructure, rows int) (data []db.FieldData, err error) {
	isNull := fieldStructure.IsNull

	switch fieldStructure.Type {
	case FieldType_INT:
		data = getTypeValue[int](fd.GetDataArray().Array, isNull)

	case FieldType_STRING:
		data = getTypeValue[string](fd.GetDataArray().Array, isNull)

	case FieldType_BOOL:
		data = getTypeValue[bool](fd.GetDataArray().Array, isNull)

	case FieldType_DATE:
		data = getTypeValue[time.Time](fd.GetDataArray().Array, isNull)

	case FieldType_REFERENCE:
		reference := fd.GetReference()
		dataArray := make([][]db.FieldData, rows)

		for i, field := range reference.GetFields() {
			fieldValues, errConv := field.ConvertToDataValues(fieldStructure.Reference.Fields[i], rows)
			if errConv != nil {
				return data, fmt.Errorf("Failed to get field values, with error: %v", err)
			}

			for row, value := range fieldValues {
				dataArray[row] = append(dataArray[row], value)
			}
		}

		data = make([]db.FieldData, rows)
		for row := range rows {
			data[row] = dataArray[row]
		}

	default:
		err = fmt.Errorf("Invalid field type of %d", fieldStructure.Type)
	}

	return data, err
}

func (fs *FieldStructure) ConvertToFieldStructure(cache *CacheStructures) (field db.Field, err error) {
	switch fs.Type {
	case FieldType_INT, FieldType_STRING, FieldType_BOOL, FieldType_DATE:
		if fieldType, err := fs.Type.ConvertToDatabase(); err == nil {
			field = db.NewPrimitiveField(fs.Name, fieldType, fs.IsNull, fs.IsKey)
		}

	case FieldType_REFERENCE:
		var refStructure *db.TableStructure
		if refStructure, err = fs.Reference.ConvertToTableStructure(cache); err == nil {
			field = db.NewReferencesField(fs.Name, refStructure, fs.IsNull, fs.IsKey)
		}

	default:
		err = fmt.Errorf("Invalied field type of %d", fs.Type)
	}

	return field, err
}

func (cs *ComponentStructure) ConvertToTableStructure(cache *CacheStructures) (*db.TableStructure, error) {
	if knowStruct, ok := cache.TableStructures[cs.Name]; ok {
		return knowStruct, nil
	}

	fields := make([]db.Field, len(cs.Fields))
	var err error

	for i, field := range cs.Fields {
		if fields[i], err = field.ConvertToFieldStructure(cache); err != nil {
			return nil, err
		}
	}

	structure := db.NewTableStructure(cs.GetName(), fields)
	cache.TableStructures[cs.Name] = structure

	return structure, nil
}

func (cd *ComponentDescription) ConvertToTableData(cache *CacheStructures) (table *db.TableData, err error) {
	if cache == nil {
		cache = NewCacheData()
	}

	structure, err := cd.Structure.ConvertToTableStructure(cache)
	if err != nil {
		return table, fmt.Errorf("Failed to get structure, with error: %v", err)
	}

	dataAmount, err := cd.Amount.GetDataBaseAmount()
	if err != nil {
		return table, fmt.Errorf("Invalied amount of %d", cd.Amount)
	}

	columns, rows := len(structure.Fields), int(cd.Rows)
	data := make([]db.FieldData, columns*rows)
	for i, field := range cd.Fields.GetFields() {
		fieldValues, err := field.ConvertToDataValues(cd.Structure.Fields[i], rows)
		if err != nil {
			return table, fmt.Errorf("Failed to get field values, with error: %v", err)
		}

		for j, dataPoint := range fieldValues {
			data[i+j*columns] = dataPoint
		}
	}

	return db.NewTableData(structure, dataAmount, data), err
}

func (cp *ComponentCompositionDescription) ConvertToCompositionDescription(cache *CacheStructures) (*db.TableComposition, error) {
	if cache == nil {
		cache = NewCacheData()
	}

	entities := make([]db.TableElement, len(cp.GetEntities()))

	for i, entityDescription := range cp.GetEntities() {
		if entity, err := entityDescription.ConvertToTableElement(cache); err != nil {
			return nil, fmt.Errorf("Failed to convert entity description while converting from table composition, with error: %v", err)

		} else {
			entities[i] = entity
		}
	}

	return db.NewTableComposition(entities...), nil
}

func (ed *EntityDescription) ConvertToTableElement(cache *CacheStructures) (tableElement db.TableElement, err error) {
	if cache == nil {
		cache = NewCacheData()
	}

	switch table := ed.GetDescription().(type) {
	case *EntityDescription_Component:
		tableElement, err = table.Component.ConvertToTableData(cache)

	case *EntityDescription_Composition:
		tableElement, err = table.Composition.ConvertToCompositionDescription(cache)

	default:
		err = fmt.Errorf("The entity description is not define")
	}

	return tableElement, err
}
