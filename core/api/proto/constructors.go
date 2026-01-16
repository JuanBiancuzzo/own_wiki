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
func NewComponentDescriptionPoint(name string, fields []*FieldDescription) *ComponentDescription {
	return &ComponentDescription{
		Name:   name,
		Amount: ComponentDescription_ONE,
		Rows:   1,
		Fields: fields,
	}
}

func NewComponentDescriptionArray(name string, rows int32, fields []*FieldDescription) *ComponentDescription {
	return &ComponentDescription{
		Name:   name,
		Amount: ComponentDescription_ARRAY,
		Rows:   rows,
		Fields: fields,
	}
}

// ---+--- FieldDescription ---+---
func NewFieldPoint(name string, typeInformation *FieldTypeInformation, point *FieldDescription_FieldData) *FieldDescription {
	return &FieldDescription{
		Name:            name,
		TypeInformation: typeInformation,
		Data: &FieldDescription_Point{
			Point: point,
		},
	}
}

func NewFieldArray(name string, typeInformation *FieldTypeInformation, array []*FieldDescription_FieldData) *FieldDescription {
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
type PrimitiveType int32

const (
	PFT_INT    = PrimitiveType(PrimitiveFieldType_INT)
	PFT_STRING = PrimitiveType(PrimitiveFieldType_STRING)
	PFT_CHAR   = PrimitiveType(PrimitiveFieldType_CHAR)
	PFT_BOOL   = PrimitiveType(PrimitiveFieldType_BOOL)
	PFT_DATE   = PrimitiveType(PrimitiveFieldType_DATE)
)

func NewPrimitiveInfo(primitiveType PrimitiveType, isNullable bool) *FieldTypeInformation {
	typeInfo := PrimitiveFieldType(primitiveType)
	if isNullable {
		typeInfo += PrimitiveFieldType_NULL_INT
	}

	return &FieldTypeInformation{
		Type: FieldTypeInformation_PRIMITIVE,
		Information: &FieldTypeInformation_Primitive{
			Primitive: typeInfo,
		},
	}
}

func NewReferenceInfo(tableName string, isNullable bool) *FieldTypeInformation {
	return &FieldTypeInformation{
		Type: FieldTypeInformation_REFERENCE,
		Information: &FieldTypeInformation_Reference{
			Reference: &ReferenceInformation{
				TableName:  tableName,
				IsNullable: isNullable,
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

func NewFieldConcreteReference(refComponent *ComponentDescription) *FieldDescription_FieldData {
	return newFieldConcreteData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Reference{
			Reference: refComponent,
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

func NewFieldNullableReference(refComponent *ComponentDescription) *FieldDescription_FieldData {
	if refComponent == nil {
		return newFieldNullableData(nil)
	}

	return newFieldNullableData(&FieldDescription_ConcreteFieldData{
		Data: &FieldDescription_ConcreteFieldData_Reference{
			Reference: refComponent,
		},
	})

}

// ---+--- ComponentCompositionDescription ---+---
func NewComponentCompositionDescription(entities []*EntityDescription) *ComponentCompositionDescription {
	return &ComponentCompositionDescription{
		Entities: entities,
	}
}
