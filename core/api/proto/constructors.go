package core

import (
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
func NewComponentDescriptionPoint(structure *ComponentStructure, data ...*FieldData) *ComponentDescription {
	fields := make([]*FieldDescription, len(data))
	for i, dataPoint := range data {
		fields[i] = NewFieldPoint(dataPoint)
	}

	return &ComponentDescription{
		Structure: structure,
		Amount:    ComponentDescription_POINT,
		Rows:      1,
		Fields: &FieldsDescription{
			Fields: fields,
		},
	}
}

func NewComponentDescriptionArray(structure *ComponentStructure, data ...[]*FieldData) *ComponentDescription {
	fields := make([]*FieldDescription, len(data))
	rows := len(data[0])
	for i, dataArray := range data {
		if len(dataArray) < rows {
			rows = len(dataArray)
		}

		fields[i] = NewFieldArray(dataArray)
	}

	return &ComponentDescription{
		Structure: structure,
		Amount:    ComponentDescription_ARRAY,
		Rows:      int32(rows),
		Fields: &FieldsDescription{
			Fields: fields,
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

// ---+--- FieldDescription ---+---
func NewFieldPoint(point *FieldData) *FieldDescription {
	return &FieldDescription{
		Data: &FieldDescription_Point{
			Point: point,
		},
	}
}

func NewFieldArray(array []*FieldData) *FieldDescription {
	return &FieldDescription{
		Data: &FieldDescription_Array{
			Array: &FieldDescription_FieldDataArray{
				DataArray: array,
			},
		},
	}
}

// ---+--- FieldData ---+---
type FnNewConcreteData[T any] func(T) *FieldData
type FnNewConcreteArray[T any] func([]T) *FieldData
type FnNewNullableData[T any] func(*T) *FieldData
type FnNewNullableArray[T any] func([]*T) *FieldData

func newConcreteData(data isFieldData_ConcreteFieldData_Data) *FieldData {
	return &FieldData{
		Data: &FieldData_Concrete{
			Concrete: &FieldData_ConcreteFieldData{
				Data: data,
			},
		},
	}
}

func newNullableData(data isFieldData_ConcreteFieldData_Data) *FieldData {
	return &FieldData{
		Data: &FieldData_Nullable{
			Nullable: &FieldData_NullableFieldData{
				Data: &FieldData_ConcreteFieldData{
					Data: data,
				},
			},
		},
	}
}

type valueConstraint interface {
	int | string | bool | time.Time |
		[]*FieldData |
		*int | *string | *bool | *time.Time
}

func NewFieldData[T valueConstraint](x T) *FieldData {
	switch value := any(x).(type) {
	case int:
		return newConcreteData(&FieldData_ConcreteFieldData_Number{Number: int32(value)})
	case *int:
		if value == nil {
			return newNullableData(nil)
		}
		return newNullableData(&FieldData_ConcreteFieldData_Number{Number: int32(*value)})

	case string:
		return newConcreteData(&FieldData_ConcreteFieldData_Text{Text: value})
	case *string:
		if value == nil {
			return newNullableData(nil)
		}
		return newNullableData(&FieldData_ConcreteFieldData_Text{Text: *value})

	case bool:
		return newConcreteData(&FieldData_ConcreteFieldData_Boolean{Boolean: value})
	case *bool:
		if value == nil {
			return newNullableData(nil)
		}
		return newNullableData(&FieldData_ConcreteFieldData_Boolean{Boolean: *value})

	case time.Time:
		return newConcreteData(&FieldData_ConcreteFieldData_Date{Date: uint32(value.Unix())})
	case *time.Time:
		if value == nil {
			return newNullableData(nil)
		}
		return newNullableData(&FieldData_ConcreteFieldData_Date{Date: uint32(value.Unix())})

	case []*FieldDescription:
		return newConcreteData(&FieldData_ConcreteFieldData_Reference{
			Reference: &FieldsDescription{
				Fields: value,
			},
		})
	case *[]*FieldDescription:
		if value == nil {
			return newNullableData(nil)
		}
		return newNullableData(&FieldData_ConcreteFieldData_Reference{
			Reference: &FieldsDescription{
				Fields: *value,
			},
		})

	default:
		panic("This is not possible")
	}
}

func NewFieldDataArray[T valueConstraint](x ...T) []*FieldData {
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
