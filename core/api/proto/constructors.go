package core

import "time"

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
func NewComponentDescriptionPoint(structure *ComponentStructure, fields []*FieldDescription) *ComponentDescription {
	return &ComponentDescription{
		Structure: structure,
		Amount:    ComponentDescription_POINT,
		Rows:      1,
		Fields: &FieldsDescription{
			Fields: fields,
		},
	}
}

func NewComponentDescriptionArray(structure *ComponentStructure, rows int32, fields []*FieldDescription) *ComponentDescription {
	return &ComponentDescription{
		Structure: structure,
		Amount:    ComponentDescription_ARRAY,
		Rows:      rows,
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

func NewReferenceStructure(name string, reference *FieldStructure, isNull, isKey bool) *FieldStructure {
	return &FieldStructure{
		Name:      name,
		Type:      FieldType_REFERENCE,
		Reference: reference,
		IsNull:    isNull,
		IsKey:     isKey,
	}
}

// ---+--- FieldDescription ---+---
func NewFieldPoint(point *FieldDescription_FieldData) *FieldDescription {
	return &FieldDescription{
		Data: &FieldDescription_Point{
			Point: point,
		},
	}
}

func NewFieldArray(array []*FieldDescription_FieldData) *FieldDescription {
	return &FieldDescription{
		Data: &FieldDescription_Array{
			Array: &FieldDescription_FieldDataArray{
				DataArray: array,
			},
		},
	}
}

// ---+--- FieldDescription_FieldData ---+---
type FnNewConcreteData[T any] func(T) *FieldDescription_FieldData
type FnNewNullableData[T any] func(*T) *FieldDescription_FieldData

func newFieldConcreteData(concrete *FieldDescription_ConcreteFieldData) *FieldDescription_FieldData {
	return &FieldDescription_FieldData{
		Data: &FieldDescription_FieldData_Concrete{
			Concrete: concrete,
		},
	}
}

func newFieldNullableData(nullable *FieldDescription_ConcreteFieldData) *FieldDescription_FieldData {
	return &FieldDescription_FieldData{
		Data: &FieldDescription_FieldData_Nullable{
			Nullable: &FieldDescription_NullableFieldData{
				Data: nullable,
			},
		},
	}
}

func NewFieldConcreteNumber(number int) *FieldDescription_FieldData {
	return newFieldConcreteData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Number{Number: int32(number)},
	})
}

func NewFieldConcreteText(text string) *FieldDescription_FieldData {
	return newFieldConcreteData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Text{Text: text},
	})
}

func NewFieldConcreteCharacter(character string) *FieldDescription_FieldData {
	return newFieldConcreteData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Character{Character: character},
	})
}

func NewFieldConcreteBoolean(boolean bool) *FieldDescription_FieldData {
	return newFieldConcreteData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Boolean{Boolean: boolean},
	})
}

func NewFieldConcreteDate(date time.Time) *FieldDescription_FieldData {
	return newFieldConcreteData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Date{Date: uint32(date.Unix())},
	})
}

func NewFieldConcreteReference(fields []*FieldDescription) *FieldDescription_FieldData {
	return newFieldConcreteData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Reference{
			Reference: &FieldsDescription{
				Fields: fields,
			},
		},
	})
}

func NewFieldNullableNumber(number *int) *FieldDescription_FieldData {
	if number == nil {
		return newFieldNullableData(nil)
	}

	return newFieldNullableData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Number{Number: int32(*number)},
	})
}

func NewFieldNullableText(text *string) *FieldDescription_FieldData {
	if text == nil {
		return newFieldNullableData(nil)
	}

	return newFieldNullableData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Text{Text: *text},
	})
}

func NewFieldNullableCharacter(character *string) *FieldDescription_FieldData {
	if character == nil {
		return newFieldNullableData(nil)
	}

	return newFieldNullableData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Character{Character: *character},
	})
}

func NewFieldNullableBoolean(boolean *bool) *FieldDescription_FieldData {
	if boolean == nil {
		return newFieldNullableData(nil)
	}

	return newFieldNullableData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Boolean{Boolean: *boolean},
	})
}

func NewFieldNullableDate(date *time.Time) *FieldDescription_FieldData {
	if date == nil {
		return newFieldNullableData(nil)
	}

	return newFieldNullableData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Date{Date: uint32(date.Unix())},
	})
}

func NewFieldNullableReference(fields []*FieldDescription) *FieldDescription_FieldData {
	if len(fields) == 0 {
		return newFieldNullableData(nil)
	}

	return newFieldNullableData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Reference{
			Reference: &FieldsDescription{
				Fields: fields,
			},
		},
	})
}

// ---+--- ComponentCompositionDescription ---+---
func NewComponentCompositionDescription(entities []*EntityDescription) *ComponentCompositionDescription {
	return &ComponentCompositionDescription{
		Entities: entities,
	}
}
