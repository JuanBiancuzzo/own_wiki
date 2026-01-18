package core

import (
	"fmt"
	"time"
)

// ---+--- From EntityDescription ---+---
func NewEntityDescriptionComponent(component *ComponentDescription) *EntityDescription {
	return &EntityDescription{
		Description: &EntityDescription_Component{
			Component: component,
		},
	}
}

func NewEntityDescriptionComposition(composition *ComponentCompositionDescription) *EntityDescription {
	return &EntityDescription{
		Description: &EntityDescription_Composition{
			Composition: composition,
		},
	}
}

// ---+--- ComponentDescription ---+---
func NewComponentDescriptionPoint(structure *ComponentStructure, fields ...*FieldDescription) *ComponentDescription {
	return &ComponentDescription{
		Structure: structure,
		Amount:    DataAmount_POINT,
		Rows:      1,
		Fields: &FieldsDescription{
			Fields: fields,
		},
	}
}

func NewComponentDescriptionArray(structure *ComponentStructure, rows int, fields ...*FieldDescription) *ComponentDescription {
	return &ComponentDescription{
		Structure: structure,
		Amount:    DataAmount_ARRAY,
		Rows:      int32(rows),
		Fields: &FieldsDescription{
			Fields: fields,
		},
	}
}

// ---+--- FieldDescription ---+---
func NewPrimitiveDescription(array ...*FieldData) *FieldDescription {
	return &FieldDescription{
		Data: &FieldDescription_DataArray{
			DataArray: &FieldDataArray{
				Array: array,
			},
		},
	}
}

func NewReferenceDescription(fields ...*FieldDescription) *FieldDescription {
	return &FieldDescription{
		Data: &FieldDescription_Reference{
			Reference: &FieldsDescription{
				Fields: fields,
			},
		},
	}
}

// ---+--- ComponentStructure ---+---
func NewComponentStructure(name string, fields []*FieldStructure) *ComponentStructure {
	return &ComponentStructure{
		Name:   name,
		Fields: fields,
	}
}

// ---+--- FieldStructure ---+---
func NewPrimitiveStructure(name string, fieldType FieldType, isNull, isKey bool) *FieldStructure {
	return &FieldStructure{
		Name:      name,
		Type:      fieldType,
		Reference: nil,
		IsNull:    isNull,
		IsKey:     isKey,
	}
}

func NewReferenceStructure(name string, reference *ComponentStructure, isNull, isKey bool) *FieldStructure {
	return &FieldStructure{
		Name:      name,
		Type:      FieldType_REFERENCE,
		Reference: reference,
		IsNull:    isNull,
		IsKey:     isKey,
	}
}

// ---+--- FieldData ---+---
func newConcreteData(data isFieldData_Data) *FieldData {
	return &FieldData{
		Data:   data,
		IsNull: false,
	}
}

func newNullableData(data isFieldData_Data, isNull bool) *FieldData {
	return &FieldData{
		Data:   data,
		IsNull: isNull,
	}
}

type valueContraint interface {
	int | string | bool | time.Time |
		*int | *string | *bool | *time.Time
}

func NewFieldData[T valueContraint](x T) *FieldData {
	switch value := any(x).(type) {
	case int:
		return newConcreteData(&FieldData_Number{Number: int32(value)})
	case *int:
		if value == nil {
			return newNullableData(&FieldData_Number{}, true)
		}
		return newNullableData(&FieldData_Number{Number: int32(*value)}, false)

	case string:
		return newConcreteData(&FieldData_Text{Text: value})
	case *string:
		if value == nil {
			return newNullableData(&FieldData_Text{}, true)
		}
		return newNullableData(&FieldData_Text{Text: *value}, false)

	case bool:
		return newConcreteData(&FieldData_Boolean{Boolean: value})
	case *bool:
		if value == nil {
			return newNullableData(&FieldData_Boolean{}, true)
		}
		return newNullableData(&FieldData_Boolean{Boolean: *value}, false)

	case time.Time:
		return newConcreteData(&FieldData_Date{Date: uint32(value.Unix())})
	case *time.Time:
		if value == nil {
			return newNullableData(&FieldData_Date{}, true)
		}
		return newNullableData(&FieldData_Date{Date: uint32(value.Unix())}, false)

	default:
		panic("This is not possible")
	}
}

func castArray[T any](array []any) ([]T, error) {
	valueArray := make([]T, len(array))
	for i, data := range array {
		if value, ok := data.(T); ok {
			valueArray[i] = value

		} else {
			return valueArray, fmt.Errorf("Mix type, expected: %T and got %T", value, data)
		}
	}
	return valueArray, nil
}

func NewFieldDataArrayWithError(x []any) (fieldData []*FieldData, err error) {
	if len(x) == 0 {
		return fieldData, fmt.Errorf("Array of length 0")
	}

	switch any(x[0]).(type) {
	case int:
		var array []int
		if array, err = castArray[int](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}
	case *int:
		var array []*int
		if array, err = castArray[*int](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}

	case string:
		var array []string
		if array, err = castArray[string](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}
	case *string:
		var array []*string
		if array, err = castArray[*string](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}

	case bool:
		var array []bool
		if array, err = castArray[bool](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}
	case *bool:
		var array []*bool
		if array, err = castArray[*bool](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}

	case time.Time:
		var array []time.Time
		if array, err = castArray[time.Time](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}
	case *time.Time:
		var array []*time.Time
		if array, err = castArray[*time.Time](x); err == nil {
			fieldData = NewFieldDataArray(array...)
		}

	default:
		err = fmt.Errorf("The type does not satisfy the concreteValueContraint or nullableValueContraint, it has type of: %T", x[0])
	}

	return fieldData, err
}

func NewFieldDataArray[T valueContraint](x ...T) []*FieldData {
	data := make([]*FieldData, len(x))
	for i, dataPoint := range x {
		data[i] = NewFieldData(dataPoint)
	}
	return data
}

// ---+--- ComponentCompositionDescription ---+---
func NewComponentCompositionDescription(entities []*EntityDescription) *ComponentCompositionDescription {
	return &ComponentCompositionDescription{
		Entities: entities,
	}
}
