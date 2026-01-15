package api_test

import (
	"testing"

	"github.com/go-test/deep"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

// ---+--- helping function ---+---

func createNullableData(concrete *pb.FieldDescription_ConcreteFieldData) *pb.FieldDescription_Point {
	return &pb.FieldDescription_Point{
		Point: &pb.FieldDescription_FieldData{
			Data: &pb.FieldDescription_FieldData_Nullable{
				Nullable: &pb.FieldDescription_NullableFieldData{
					Data: concrete,
				},
			},
		},
	}
}

func createConcretePointData(concrete *pb.FieldDescription_ConcreteFieldData) *pb.FieldDescription_Point {
	return &pb.FieldDescription_Point{
		Point: &pb.FieldDescription_FieldData{
			Data: &pb.FieldDescription_FieldData_Concrete{
				Concrete: concrete,
			},
		},
	}
}

func createConcreteArrayData(concreteArray []*pb.FieldDescription_ConcreteFieldData) *pb.FieldDescription_Array {
	dataArray := make([]*pb.FieldDescription_FieldData, len(concreteArray))
	for i, concrete := range concreteArray {
		dataArray[i] = &pb.FieldDescription_FieldData{
			Data: &pb.FieldDescription_FieldData_Concrete{
				Concrete: concrete,
			},
		}
	}

	return &pb.FieldDescription_Array{
		Array: &pb.FieldDescription_FieldDataArray{
			DataArray: dataArray,
		},
	}
}

func createNullableArrayData(concreteArray []*pb.FieldDescription_ConcreteFieldData) *pb.FieldDescription_Array {
	dataArray := make([]*pb.FieldDescription_FieldData, len(concreteArray))
	for i, concrete := range concreteArray {
		dataArray[i] = &pb.FieldDescription_FieldData{
			Data: &pb.FieldDescription_FieldData_Nullable{
				Nullable: &pb.FieldDescription_NullableFieldData{
					Data: concrete,
				},
			},
		}
	}

	return &pb.FieldDescription_Array{
		Array: &pb.FieldDescription_FieldDataArray{
			DataArray: dataArray,
		},
	}
}

// ---+--- tests ---+---

func TestConventSimpelTableDataWithSinglePointAndPrimitiveValues(t *testing.T) {
	/* This is the representation of the structure:
	type TestTable struct {
		Number  int
		String  string
		Boolean bool
	}
	*/
	var tableElement db.TableElement = db.TableData{
		Structure: db.TableStructure{
			Name: "TestTable",
			Fields: []db.Field{
				db.NewPrimitiveField("Number", db.FT_INT, false, false),
				db.NewPrimitiveField("String", db.FT_STRING, false, false),
				db.NewPrimitiveField("Boolean", db.FT_BOOL, false, false),
			},
		},
		DataAmount: db.DA_POINT,
		Data:       []db.FieldData{int(1), "Text", true},
	}

	entityDescriptionExpected := &pb.EntityDescription{
		Description: &pb.EntityDescription_Component{
			Component: &pb.ComponentDescription{
				Name:   "TestTable",
				Amount: pb.ComponentDescription_ONE,
				Rows:   1,
				Fields: []*pb.FieldDescription{
					{
						Name: "Number",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_INT,
							},
						},
						Data: createConcretePointData(&pb.FieldDescription_ConcreteFieldData{
							Data: &pb.FieldDescription_ConcreteFieldData_Number{Number: 1},
						}),
					},
					{
						Name: "String",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_STRING,
							},
						},
						Data: createConcretePointData(&pb.FieldDescription_ConcreteFieldData{
							Data: &pb.FieldDescription_ConcreteFieldData_Text{Text: "Text"},
						}),
					},
					{
						Name: "Boolean",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_BOOL,
							},
						},
						Data: createConcretePointData(&pb.FieldDescription_ConcreteFieldData{
							Data: &pb.FieldDescription_ConcreteFieldData_Boolean{Boolean: true},
						}),
					},
				},
			},
		},
	}

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(); err != nil {
		t.Errorf("While converting to TableElement, got the error: %v", err)

	} else if diff := deep.Equal(tableElement, tableElementGen); diff != nil {
		t.Error(diff)
	}
}

func TestConventSimpelTableDataWithSinglePointAndPrimitiveValuesAndNullable(t *testing.T) {
	/* This is the representation of the structure:
	type TestTable struct {
		Number  Nullable[int]
		String  string
		Date    Nullable[uint]
	}
	*/
	var numberValue int = 1
	var dateValue *uint = nil

	var tableElement db.TableElement = db.TableData{
		Structure: db.TableStructure{
			Name: "TestTable",
			Fields: []db.Field{
				db.NewPrimitiveField("Number", db.FT_INT, true, false),
				db.NewPrimitiveField("String", db.FT_STRING, false, false),
				db.NewPrimitiveField("Date", db.FT_DATE, true, false),
			},
		},
		DataAmount: db.DA_POINT,
		Data:       []db.FieldData{&numberValue, "Text", dateValue},
	}

	entityDescriptionExpected := &pb.EntityDescription{
		Description: &pb.EntityDescription_Component{
			Component: &pb.ComponentDescription{
				Name:   "TestTable",
				Amount: pb.ComponentDescription_ONE,
				Rows:   1,
				Fields: []*pb.FieldDescription{
					{
						Name: "Number",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_NULL_INT,
							},
						},
						Data: createNullableData(&pb.FieldDescription_ConcreteFieldData{
							Data: &pb.FieldDescription_ConcreteFieldData_Number{Number: 1},
						}),
					},
					{
						Name: "String",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_STRING,
							},
						},
						Data: createConcretePointData(&pb.FieldDescription_ConcreteFieldData{
							Data: &pb.FieldDescription_ConcreteFieldData_Text{Text: "Text"},
						}),
					},
					{
						Name: "Date",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_NULL_DATE,
							},
						},
						Data: createNullableData(nil),
					},
				},
			},
		},
	}

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(); err != nil {
		t.Errorf("While converting to TableElement, got the error: %v", err)

	} else if diff := deep.Equal(tableElement, tableElementGen); diff != nil {
		t.Error(diff)
	}
}

func TestConventSimpelTableDataWithArrayAndPrimitiveValues(t *testing.T) {
	/* This is the representation an array of the structure:
	type TestTable struct {
		Number  Nullable[int]
		String  string
		Boolean bool
	}
	*/
	numbers := []int{1, 2, 3}
	var nullNumber *int = nil

	var tableElement db.TableElement = db.TableData{
		Structure: db.TableStructure{
			Name: "TestTable",
			Fields: []db.Field{
				db.NewPrimitiveField("Number", db.FT_INT, true, false),
				db.NewPrimitiveField("String", db.FT_STRING, false, false),
				db.NewPrimitiveField("Boolean", db.FT_BOOL, false, false),
			},
		},
		DataAmount: db.DA_ARRAY,
		Data: []db.FieldData{
			&numbers[0], "Primero", true,
			nullNumber, "Segundo", true,
			&numbers[2], "Tercero", false,
		},
	}

	entityDescriptionExpected := &pb.EntityDescription{
		Description: &pb.EntityDescription_Component{
			Component: &pb.ComponentDescription{
				Name:   "TestTable",
				Amount: pb.ComponentDescription_ARRAY,
				Rows:   3,
				Fields: []*pb.FieldDescription{
					{
						Name: "Number",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_NULL_INT,
							},
						},
						Data: createNullableArrayData([]*pb.FieldDescription_ConcreteFieldData{
							{Data: &pb.FieldDescription_ConcreteFieldData_Number{Number: 1}},
							nil,
							{Data: &pb.FieldDescription_ConcreteFieldData_Number{Number: 3}},
						}),
					},
					{
						Name: "String",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_STRING,
							},
						},
						Data: createConcreteArrayData([]*pb.FieldDescription_ConcreteFieldData{
							{Data: &pb.FieldDescription_ConcreteFieldData_Text{Text: "Primero"}},
							{Data: &pb.FieldDescription_ConcreteFieldData_Text{Text: "Segundo"}},
							{Data: &pb.FieldDescription_ConcreteFieldData_Text{Text: "Tercero"}},
						}),
					},
					{
						Name: "Boolean",
						TypeInformation: &pb.FieldTypeInformation{
							Type: pb.FieldTypeInformation_PRIMITIVE,
							Information: &pb.FieldTypeInformation_Primitive{
								Primitive: pb.PrimitiveFieldType_BOOL,
							},
						},
						Data: createConcreteArrayData([]*pb.FieldDescription_ConcreteFieldData{
							{Data: &pb.FieldDescription_ConcreteFieldData_Boolean{Boolean: true}},
							{Data: &pb.FieldDescription_ConcreteFieldData_Boolean{Boolean: true}},
							{Data: &pb.FieldDescription_ConcreteFieldData_Boolean{Boolean: false}},
						}),
					},
				},
			},
		},
	}

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(); err != nil {
		t.Errorf("While converting to TableElement, got the error: %v", err)

	} else if diff := deep.Equal(tableElement, tableElementGen); diff != nil {
		t.Error(diff)
	}
}
