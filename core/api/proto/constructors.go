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
func NewComponentDescription(name string, amount ComponentDescription_DataAmount, rows int32, fields []*FieldDescription) *ComponentDescription {
	return &ComponentDescription{
		Name:   name,
		Amount: amount,
		Rows:   rows,
		Fields: fields,
	}
}

// ---+--- FieldDescription ---+---
func NewFieldDescriptionPoint(name string, typeInformation *FieldTypeInformation, point *FieldDescription_FieldData) *FieldDescription {
	return &FieldDescription{
		Name:            name,
		TypeInformation: typeInformation,
		Data: &FieldDescription_Point{
			Point: point,
		},
	}
}

func NewFieldDescriptionArray(name string, typeInformation *FieldTypeInformation, array []*FieldDescription_FieldData) *FieldDescription {
	return &FieldDescription{
		Name:            name,
		TypeInformation: typeInformation,
		Data: &FieldDescription_Array{
			Array: &FieldDescription_FieldDataArray{
				DataArray: array,
			},
		},
	}
}

// ---+--- TypeInformation ---+---
func NewFieldTypeInformationPrimitive(primitiveType PrimitiveFieldType) *FieldTypeInformation {
	return &FieldTypeInformation{
		Type: FieldTypeInformation_PRIMITIVE,
		Information: &FieldTypeInformation_Primitive{
			Primitive: primitiveType,
		},
	}
}

func NewFieldTypeInformationReference(tableName string) *FieldTypeInformation {
	return &FieldTypeInformation{
		Type: FieldTypeInformation_REFERENCE,
		Information: &FieldTypeInformation_Reference{
			Reference: &ReferenceInformation{
				TableName: tableName,
			},
		},
	}
}

// ---+--- FieldDescription_FieldData ---+---
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

// ---+--- ComponentCompositionDescription ---+---
func NewComponentCompositionDescription(entities []*EntityDescription) *ComponentCompositionDescription {
	return &ComponentCompositionDescription{
		Entities: entities,
	}
}
