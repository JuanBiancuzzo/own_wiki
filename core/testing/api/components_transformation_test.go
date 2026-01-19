package api_test

import (
	"testing"
	"time"

	"github.com/go-test/deep"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	db "github.com/JuanBiancuzzo/own_wiki/core/database"
)

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
		pb.NewComponentDescriptionPoint(
			pb.NewComponentStructure("TestTable", []*pb.FieldStructure{
				pb.NewPrimitiveStructure("Number", pb.FieldType_INT, false, false),
				pb.NewPrimitiveStructure("String", pb.FieldType_STRING, false, false),
				pb.NewPrimitiveStructure("Boolean", pb.FieldType_BOOL, false, false),
			}),
			pb.NewPrimitiveDescription(pb.NewFieldData(1)),
			pb.NewPrimitiveDescription(pb.NewFieldData("Text")),
			pb.NewPrimitiveDescription(pb.NewFieldData(true)),
		),
	)

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(nil); err != nil {
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
	var numberValue *int = nil
	var dateValue *time.Time = nil

	var tableElement db.TableElement = db.NewTableData(
		db.NewTableStructure("TestTable", []db.Field{
			db.NewPrimitiveField("Number", db.FT_INT, true, false),
			db.NewPrimitiveField("String", db.FT_STRING, false, false),
			db.NewPrimitiveField("Date", db.FT_DATE, true, false),
		}),
		db.DA_POINT,
		[]db.FieldData{numberValue, "Text", dateValue},
	)

	entityDescriptionExpected := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionPoint(
			pb.NewComponentStructure("TestTable", []*pb.FieldStructure{
				pb.NewPrimitiveStructure("Number", pb.FieldType_INT, true, false),
				pb.NewPrimitiveStructure("String", pb.FieldType_STRING, false, false),
				pb.NewPrimitiveStructure("Date", pb.FieldType_DATE, true, false),
			}),
			pb.NewPrimitiveDescription(pb.NewFieldData(numberValue)),
			pb.NewPrimitiveDescription(pb.NewFieldData("Text")),
			pb.NewPrimitiveDescription(pb.NewFieldData(dateValue)),
		),
	)

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(nil); err != nil {
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
		pb.NewComponentDescriptionArray(
			pb.NewComponentStructure("TestTable", []*pb.FieldStructure{
				pb.NewPrimitiveStructure("Number", pb.FieldType_INT, true, false),
				pb.NewPrimitiveStructure("String", pb.FieldType_STRING, false, false),
				pb.NewPrimitiveStructure("Boolean", pb.FieldType_BOOL, false, false),
			}),
			3,
			pb.NewPrimitiveDescription(pb.NewFieldDataArray(&numbers[0], nullNumber, &numbers[2])...),
			pb.NewPrimitiveDescription(pb.NewFieldDataArray("Primero", "Segundo", "Tercero")...),
			pb.NewPrimitiveDescription(pb.NewFieldDataArray(true, true, false)...),
		),
	)

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(nil); err != nil {
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

	referenceTableStructure := db.NewTableStructure("ReferencesTable", []db.Field{
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

	refStructure := pb.NewComponentStructure("ReferencesTable", []*pb.FieldStructure{
		pb.NewPrimitiveStructure("Some", pb.FieldType_INT, false, false),
		pb.NewPrimitiveStructure("Thing", pb.FieldType_INT, false, false),
	})

	entityDescriptionExpected := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionPoint(
			pb.NewComponentStructure("TestTable", []*pb.FieldStructure{
				pb.NewPrimitiveStructure("Number", pb.FieldType_INT, false, false),
				pb.NewPrimitiveStructure("String", pb.FieldType_STRING, false, false),
				pb.NewReferenceStructure("Reference", refStructure, false, false),
			}),
			pb.NewPrimitiveDescription(pb.NewFieldData(1)),
			pb.NewPrimitiveDescription(pb.NewFieldData("Text")),
			pb.NewReferenceDescription(
				pb.NewPrimitiveDescription(pb.NewFieldData(2)),
				pb.NewPrimitiveDescription(pb.NewFieldData(3)),
			),
		),
	)

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(nil); err != nil {
		t.Errorf("While converting to TableElement, got the error: %v", err)

	} else if diff := deep.Equal(tableElement, tableElementGen); diff != nil {
		t.Error(diff)
	}
}

func TestConventCompositeTableDataWithArrayPointAndAllValues(t *testing.T) {
	/* This is the representation of the structure:
	Type TestCompositeTable struct {
		TestTable
		[]OtherTable
	}

	type TestTable struct {
		Number    int
		String    string
		Reference *ReferencesTable
	}

	type OtherTable struct {
		Number  Optional[int]
		Boolean bool
	}

	type ReferencesTable struct {
		Some int
		Thing int
	}
	*/

	numbers := []int{4, 5, 6}
	var nilNumber *int = nil

	otherTable := db.NewTableData(
		db.NewTableStructure("OtherTable", []db.Field{
			db.NewPrimitiveField("Number", db.FT_INT, true, false),
			db.NewPrimitiveField("Boolean", db.FT_BOOL, false, false),
		}),
		db.DA_ARRAY,
		[]db.FieldData{
			&numbers[0], true,
			nilNumber, false,
			&numbers[1], false,
			&numbers[2], true,
		},
	)

	referenceTableStructure := db.NewTableStructure("ReferencesTable", []db.Field{
		db.NewPrimitiveField("Some", db.FT_INT, false, false),
		db.NewPrimitiveField("Thing", db.FT_INT, false, false),
	})

	testTable := db.NewTableData(
		db.NewTableStructure("TestTable", []db.Field{
			db.NewPrimitiveField("Number", db.FT_INT, false, false),
			db.NewPrimitiveField("String", db.FT_STRING, false, false),
			db.NewReferencesField("Reference", referenceTableStructure, false, false),
		}),
		db.DA_POINT,
		[]db.FieldData{1, "Text", []db.FieldData{2, 3}},
	)

	var tableElement db.TableElement = db.NewTableComposition(testTable, otherTable)

	testOtherDescription := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionArray(
			pb.NewComponentStructure("OtherTable", []*pb.FieldStructure{
				pb.NewPrimitiveStructure("Number", pb.FieldType_INT, true, false),
				pb.NewPrimitiveStructure("Boolean", pb.FieldType_BOOL, false, false),
			}),
			4,
			pb.NewPrimitiveDescription(pb.NewFieldDataArray(&numbers[0], nilNumber, &numbers[1], &numbers[2])...),
			pb.NewPrimitiveDescription(pb.NewFieldDataArray(true, false, false, true)...),
		),
	)

	refStructure := pb.NewComponentStructure("ReferencesTable", []*pb.FieldStructure{
		pb.NewPrimitiveStructure("Some", pb.FieldType_INT, false, false),
		pb.NewPrimitiveStructure("Thing", pb.FieldType_INT, false, false),
	})

	testEntityDescription := pb.NewEntityDescriptionComponent(
		pb.NewComponentDescriptionPoint(
			pb.NewComponentStructure("TestTable", []*pb.FieldStructure{
				pb.NewPrimitiveStructure("Number", pb.FieldType_INT, false, false),
				pb.NewPrimitiveStructure("String", pb.FieldType_STRING, false, false),
				pb.NewReferenceStructure("Reference", refStructure, false, false),
			}),
			pb.NewPrimitiveDescription(pb.NewFieldData(1)),
			pb.NewPrimitiveDescription(pb.NewFieldData("Text")),
			pb.NewReferenceDescription(
				pb.NewPrimitiveDescription(pb.NewFieldData(2)),
				pb.NewPrimitiveDescription(pb.NewFieldData(3)),
			),
		),
	)

	entityDescriptionExpected := pb.NewEntityDescriptionComposition(
		pb.NewComponentCompositionDescription(testEntityDescription, testOtherDescription),
	)

	if entityDescription, err := pb.ConvertToEntityDescription(&tableElement); err != nil {
		t.Errorf("While converting to EntityDescription, got the error: %v", err)

	} else if diff := deep.Equal(entityDescriptionExpected, entityDescription); diff != nil {
		t.Error(diff)

	} else if tableElementGen, err := entityDescription.ConvertToTableElement(nil); err != nil {
		t.Errorf("While converting to TableElement, got the error: %v", err)

	} else if diff := deep.Equal(tableElement, tableElementGen); diff != nil {
		t.Error(diff)
	}
}
