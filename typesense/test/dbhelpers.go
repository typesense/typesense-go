// +build integration

package test

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/typesense/typesense-go/typesense/api"
)

func getNewCollectionName(namePrefix string) string {
	nameUUID := uuid.New()
	return fmt.Sprintf("%s_%s", namePrefix, nameUUID.String())
}

func createNewSchema(collectionName string) *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{
				Name: "company_name",
				Type: "string",
			},
			{
				Name: "num_employees",
				Type: "int32",
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: true,
			},
		},
		DefaultSortingField: "num_employees",
	}
}

func expectedNewCollection(name string) *api.Collection {
	return &api.Collection{
		CollectionSchema: *createNewSchema(name),
		NumDocuments:     0,
	}
}
