package typesense

import (
	"errors"
	"reflect"
	"strings"

	"github.com/typesense/typesense-go/typesense/api"
)

// CollectionNamer is an interface that provides a method to get the collection name.
type CollectionNamer interface {
	CollectionName() string
}

// CreateSchemaFromGoStruct takes a Go struct and generates a Typesense CollectionSchema.
// If the struct implements the CollectionNamer interface, its CollectionName method is used to get the collection name.
func CreateSchemaFromGoStruct(structData interface{}) (*api.CollectionSchema, error) {
	t := reflect.TypeOf(structData)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var collectionName string
	if namer, ok := structData.(CollectionNamer); ok {
		collectionName = namer.CollectionName()
	} else {
		collectionName = t.Name()
	}

	fields := make([]api.Field, 0)
	var defaultSortingField *string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue, ok := field.Tag.Lookup("typesense")
		if !ok || tagValue == "-" {
			continue
		}

		fieldType := field.Type.String()
		if fieldType == "uuid.UUID" {
			fieldType = "string"
		}

		tagParts := strings.Split(tagValue, ",")
		facetValue := false // Default facet value
		typesenseField := api.Field{
			Name:  field.Name,
			Type:  fieldType,
			Facet: &facetValue, // Initially false
		}

		for _, tagPart := range tagParts {
			tagPartTrimmed := strings.TrimSpace(tagPart)
			if tagPartTrimmed == "defaultSort" {
				if defaultSortingField != nil {
					return nil, errors.New("multiple fields marked with 'defaultSort' tag")
				}
				defaultSortingField = &field.Name
			} else if tagPartTrimmed == "facet" {
				facetValue = true
				typesenseField.Facet = &facetValue
			}
		}

		fields = append(fields, typesenseField)
	}

	return &api.CollectionSchema{
		Name:                collectionName,
		Fields:              fields,
		DefaultSortingField: defaultSortingField,
	}, nil
}
