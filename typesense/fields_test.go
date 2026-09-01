package typesense

import (
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/v2/typesense/api"
	"testing"
)

type GenerationTest struct {
	ID            string `json:"id,index"`
	Name          string `json:"name,index,sort"`
	UserId        string `json:"user_id,index,join:user.id"` // creates a reference to the collection use
	Birthdate     int64  `json:"birthdate"`
	LastTreatment int64  `json:"last_treatment,index,optional"`
	LocationId    string `json:"location_id,facet"`
}

type User struct {
	ID   string `json:"id,index"`
	Name string `json:"name,index"`
	Type int32  `json:"type,facet"`
}

func TestToFields_GenerationTest(t *testing.T) {
	testStruct := GenerationTest{}
	expectedFields := []api.Field{
		{
			Name:  "id",
			Type:  "string",
			Index: Pointer(true),
		},
		{
			Name:  "name",
			Type:  "string",
			Index: Pointer(true),
		},
		{
			Name:      "user_id",
			Type:      "string",
			Index:     Pointer(true),
			Reference: Pointer("user.id"),
		},
		{
			Name: "birthdate",
			Type: "int64",
		},
		{
			Name:     "last_treatment",
			Type:     "int64",
			Index:    Pointer(true),
			Optional: Pointer(true),
		},
		{
			Name:  "location_id",
			Type:  "string",
			Facet: Pointer(true),
		},
	}

	fields, err := ToFields(testStruct)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	compareFields(t, expectedFields, fields)
}

// Helper function to compare slices of api.Field
func compareFields(t *testing.T, expected, actual []api.Field) {
	assert.Equal(t, len(expected), len(actual))

	for i, exp := range expected {
		act := actual[i]
		assert.Equal(t, exp.Type, act.Type)
		assert.Equal(t, exp.Name, act.Name)
		if exp.Index != nil && act.Index != nil {
			assert.Equal(t, *exp.Index, *act.Index)
		}
	}
}
