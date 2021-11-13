package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type yml map[string]interface{}

func main() {
	m := make(yml)

	configFile, err := os.Open("./openapi.yml")
	if err != nil {
		log.Fatalf("Unable to open config file: %s", err.Error())
		return
	}

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unwrapping the search parameters
	unwrapSearchParameters(&m)
	// Unwrapping import and export parameters
	unwrapImportDocuments(&m)
	unwrapExportDocuments(&m)
	// Unwrapping delete document parameters
	unwrapDeleteDocument(&m)
	// Remove additionalProperties from SearchResultHit -> document
	searchResultHit(&m)

	generatorFile, err := os.Create("./generator.yml")
	if err != nil {
		log.Fatalf("Unable to open config file: %s", err.Error())
		return
	}

	encode := yaml.NewEncoder(generatorFile)
	encode.SetIndent(2)
	err = encode.Encode(m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func searchResultHit(m *yml) {
	properties := (*m)["components"].(yml)["schemas"].(yml)["SearchResultHit"].(yml)["properties"].(yml)
	document := properties["document"].(yml)
	delete(document, "additionalProperties")
}

func unwrapDeleteDocument(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents"].(yml)["delete"].(yml)["parameters"].([]interface{})
	deleteParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)
	for k, v := range deleteParameters {
		newMap := make(yml)
		newMap["name"] = k
		newMap["in"] = "query"
		newMap["schema"] = make(yml)
		newMap["schema"].(yml)["type"] = v.(yml)["type"].(string)
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents"].(yml)["delete"].(yml)["parameters"] = parameters
}

func unwrapExportDocuments(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents/export"].(yml)["get"].(yml)["parameters"].([]interface{})
	exportParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)
	for k, v := range exportParameters {
		newMap := make(yml)
		newMap["name"] = k
		newMap["in"] = "query"
		newMap["schema"] = make(yml)
		if v.(yml)["type"].(string) == "array" {
			newMap["schema"].(yml)["type"] = "array"
			newMap["schema"].(yml)["items"] = v.(yml)["items"]
		} else {
			newMap["schema"].(yml)["type"] = v.(yml)["type"].(string)
		}
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents/export"].(yml)["get"].(yml)["parameters"] = parameters
}

func unwrapImportDocuments(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents/import"].(yml)["post"].(yml)["parameters"].([]interface{})
	importParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)

	for k, v := range importParameters {
		newMap := make(yml)
		newMap["name"] = k
		newMap["in"] = "query"
		newMap["schema"] = make(yml)
		newMap["schema"].(yml)["type"] = v.(yml)["type"].(string)
		if v.(yml)["enum"] != nil {
			newMap["schema"].(yml)["enum"] = v.(yml)["enum"]
		}
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents/import"].(yml)["post"].(yml)["parameters"] = parameters
}

func unwrapSearchParameters(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents/search"].(yml)["get"].(yml)["parameters"].([]interface{})
	searchParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)

	for k, v := range searchParameters {
		newMap := make(yml)
		newMap["name"] = k
		if k == "q" || k == "query_by" {
			newMap["required"] = true
		}
		newMap["in"] = "query"
		newMap["schema"] = make(yml)
		if v.(yml)["oneOf"] == nil {
			if v.(yml)["type"].(string) == "array" {
				newMap["schema"].(yml)["type"] = "array"
				newMap["schema"].(yml)["items"] = v.(yml)["items"]
			} else {
				newMap["schema"].(yml)["type"] = v.(yml)["type"].(string)
			}
		} else {
			newMap["schema"].(yml)["oneOf"] = v.(yml)["oneOf"]
		}
		parameters = append(parameters, newMap)
	}

	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents/search"].(yml)["get"].(yml)["parameters"] = parameters
}
