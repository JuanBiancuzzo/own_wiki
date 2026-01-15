package api_test

import (
	"testing"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

/*
We want to test the conversion from EntityDescription to TableElement:
  - If
*/
func TestConvertFromEntityDescriptionToTableElement(t *testing.T) {
}

/*
We want to test the conversion from TableElement to EntityDescription:
  - The element is a simple TableData with a single point, with primitive values
*/
func TestConvertFromTableElementToEntityDescription(t *testing.T) {
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

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting, got the error: %v", err)

	} else if componentDescription, ok := entityDescription.Description.(*pb.EntityDescription_Component); !ok {
		t.Error("The entity description was not a component")

	} else if component := componentDescription.Component; component.Name != "TestTable" {
		t.Errorf("The component name was not set correctly, should be 'TestTable' but was set to %s", component.Name)

	} else if component.Amount != pb.ComponentDescription_ONE {
		t.Error("The amount of element was not set correctly")

	} else if component.Rows != 1 {
		t.Errorf("The rows should be one, but was %d", component.Rows)

	} else if len(component.Fields) != 3 {
		t.Errorf("The amount of field should be 3, but was %d", len(component.Fields))

	} else if numberField := component.Fields[0]; numberField.Name != "Number" {
		t.Errorf("The name of the first field should be 'Number' but was set to %s", numberField.Name)

	} else if dataPoint, ok := numberField.Data.(*pb.FieldDescription_Point); !ok {
		t.Error("The data of the number was not a point")

	} else if concrete, ok := dataPoint.Point.Data.(*pb.FieldDescription_FieldData_Concrete); !ok {
		t.Errorf("The number should be a concrete value, but was %T", dataPoint.Point.Data)

	} else if concreteNumber, ok := concrete.Concrete.Data.(*pb.FieldDescription_ConcreteFieldData_Number); !ok {
		t.Error("the number should be a concrete number")

	} else if int(concreteNumber.Number) != 1 {
		t.Errorf("the number should should be set to 1, but was set to %d", concreteNumber.Number)

	} else if textField := component.Fields[1]; textField.Name != "String" {
		t.Errorf("The name of the first field should be 'String' but was set to %s", textField.Name)

	} else if dataPoint, ok := textField.Data.(*pb.FieldDescription_Point); !ok {
		t.Error("The data of the text was not a point")

	} else if concrete, ok := dataPoint.Point.Data.(*pb.FieldDescription_FieldData_Concrete); !ok {
		t.Errorf("The text should be a concrete value, but was %T", dataPoint.Point.Data)

	} else if concreteText, ok := concrete.Concrete.Data.(*pb.FieldDescription_ConcreteFieldData_Text); !ok {
		t.Error("the number should be a concrete text")

	} else if concreteText.Text != "Text" {
		t.Errorf("the number should should be set to 'Text', but was set to %s", concreteText.Text)

	} else if booleanField := component.Fields[2]; booleanField.Name != "Boolean" {
		t.Errorf("The name of the first field should be 'Boolean' but was set to %s", booleanField.Name)

	} else if dataPoint, ok := booleanField.Data.(*pb.FieldDescription_Point); !ok {
		t.Error("The data of the boolean was not a point")

	} else if concrete, ok := dataPoint.Point.Data.(*pb.FieldDescription_FieldData_Concrete); !ok {
		t.Errorf("The boolean should be a concrete value, but was %T", dataPoint.Point.Data)

	} else if concreteBoolean, ok := concrete.Concrete.Data.(*pb.FieldDescription_ConcreteFieldData_Boolean); !ok {
		t.Error("the boolean should be a concrete boolean")

	} else if concreteBoolean.Boolean != true {
		t.Error("the number should should be set to 'true', but was set to false")
	}
}
