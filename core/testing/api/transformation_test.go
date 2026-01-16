package api_test

import (
	"testing"
	"time"

	"github.com/go-test/deep"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

// ---+--- helping function ---+---

// ---+--- tests ---+---

func TestConventSimpelTableDataWithSinglePointAndPrimitiveValues(t *testing.T) {
	/* This is the representation of the structure:
	type TestTable struct {
		Number  int
		String  string
		Boolean bool
	}
	*/
	var tableElement db.TableElement = db.NewTableData(
		db.NewTableStructure("TestTable", []db.Field{
			db.NewPrimitiveField("Number", db.FT_INT, false, false),
			db.NewPrimitiveField("String", db.FT_STRING, false, false),
			db.NewPrimitiveField("Boolean", db.FT_BOOL, false, false),
		}),
		db.DA_POINT,
		[]db.FieldData{int(1), "Text", true},
	)

	entityDescriptionExpected := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionPoint("TestTable", []*pb.FieldDescription{
			pb.NewFieldPoint("Number", pb.NewPrimitiveInfo(pb.PFT_INT, false), pb.NewFieldConcreteNumber(1)),
			pb.NewFieldPoint("String", pb.NewPrimitiveInfo(pb.PFT_STRING, false), pb.NewFieldConcreteText("Text")),
			pb.NewFieldPoint("Boolean", pb.NewPrimitiveInfo(pb.PFT_BOOL, false), pb.NewFieldConcreteBoolean(true)),
		}),
	)

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
		Date    Nullable[time.Time]
	}
	*/
	var numberValue int = 1
	var dateValue *time.Time = nil

	var tableElement db.TableElement = db.NewTableData(
		db.NewTableStructure("TestTable", []db.Field{
			db.NewPrimitiveField("Number", db.FT_INT, true, false),
			db.NewPrimitiveField("String", db.FT_STRING, false, false),
			db.NewPrimitiveField("Date", db.FT_DATE, true, false),
		}),
		db.DA_POINT,
		[]db.FieldData{&numberValue, "Text", dateValue},
	)

	entityDescriptionExpected := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionPoint("TestTable", []*pb.FieldDescription{
			pb.NewFieldPoint("Number", pb.NewPrimitiveInfo(pb.PFT_INT, true), pb.NewFieldNullableNumber(&numberValue)),
			pb.NewFieldPoint("String", pb.NewPrimitiveInfo(pb.PFT_STRING, false), pb.NewFieldConcreteText("Text")),
			pb.NewFieldPoint("Date", pb.NewPrimitiveInfo(pb.PFT_DATE, true), pb.NewFieldNullableDate(nil)),
		}),
	)

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

	var tableElement db.TableElement = db.NewTableData(
		db.NewTableStructure("TestTable", []db.Field{
			db.NewPrimitiveField("Number", db.FT_INT, true, false),
			db.NewPrimitiveField("String", db.FT_STRING, false, false),
			db.NewPrimitiveField("Boolean", db.FT_BOOL, false, false),
		}),
		db.DA_ARRAY,
		[]db.FieldData{
			&numbers[0], "Primero", true,
			nullNumber, "Segundo", true,
			&numbers[2], "Tercero", false,
		},
	)
	entityDescriptionExpected := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionArray("TestTable", 3, []*pb.FieldDescription{
			pb.NewFieldArray("Number", pb.NewPrimitiveInfo(pb.PFT_INT, true), []*pb.FieldDescription_FieldData{
				pb.NewFieldNullableNumber(&numbers[0]), pb.NewFieldNullableNumber(nil), pb.NewFieldNullableNumber(&numbers[2]),
			}),
			pb.NewFieldArray("String", pb.NewPrimitiveInfo(pb.PFT_STRING, false), []*pb.FieldDescription_FieldData{
				pb.NewFieldConcreteText("Primero"), pb.NewFieldConcreteText("Segundo"), pb.NewFieldConcreteText("Tercero"),
			}),
			pb.NewFieldArray("Boolean", pb.NewPrimitiveInfo(pb.PFT_BOOL, false), []*pb.FieldDescription_FieldData{
				pb.NewFieldConcreteBoolean(true), pb.NewFieldConcreteBoolean(true), pb.NewFieldConcreteBoolean(false),
			}),
		}),
	)

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

func TestConventSimpelTableDataWithPointAndAllValues(t *testing.T) {
	/* This is the representation of the structure:
	type TestTable struct {
		Number    int
		String    string
		Reference *ReferencesTable
	}

	type ReferencesTable struct {
		Some int
		Thing int
	}
	*/

	referenceTableStructure := db.NewTableStructure("RefererncesTable", []db.Field{
		db.NewPrimitiveField("Some", db.FT_INT, false, false),
		db.NewPrimitiveField("Thing", db.FT_INT, false, false),
	})

	var tableElement db.TableElement = db.NewTableData(
		db.NewTableStructure("TestTable", []db.Field{
			db.NewPrimitiveField("Number", db.FT_INT, false, false),
			db.NewPrimitiveField("String", db.FT_STRING, false, false),
			db.NewReferencesField("Reference", referenceTableStructure, false, false),
		}),
		db.DA_POINT,
		[]db.FieldData{1, "Text", []db.FieldData{2, 3}},
	)

	refComponentExpected := pb.NewComponentDescriptionPoint("ReferenceTable", []*pb.FieldDescription{
		pb.NewFieldPoint("Some", pb.NewPrimitiveInfo(pb.PFT_INT, false), pb.NewFieldConcreteNumber(2)),
		pb.NewFieldPoint("Thing", pb.NewPrimitiveInfo(pb.PFT_INT, false), pb.NewFieldConcreteNumber(3)),
	})

	entityDescriptionExpected := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionPoint("TestTable", []*pb.FieldDescription{
			pb.NewFieldPoint("Number", pb.NewPrimitiveInfo(pb.PFT_INT, false), pb.NewFieldConcreteNumber(1)),
			pb.NewFieldPoint("String", pb.NewPrimitiveInfo(pb.PFT_STRING, false), pb.NewFieldConcreteText("Text")),
			pb.NewFieldPoint("Reference", pb.NewReferenceInfo("ReferencesTable", false), pb.NewFieldConcreteReference(refComponentExpected)),
		}),
	)

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
