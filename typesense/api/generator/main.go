package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type yml map[string]interface{}

const (
	query = "query"
	array = "array"
)

type MapKV struct {
	Key   string
	Value interface{}
}

func sortedSlice(params map[string]interface{}) []MapKV {
	kvs := []MapKV{}

	for k, v := range params {
		kvs = append(kvs, MapKV{k, v})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Key < kvs[j].Key
	})

	return kvs
}

// This script makes the changes needed for oapi-codegen to generate client_gen.go and types_gen.go from
// https://github.com/typesense/typesense-api-spec/blob/master/openapi.yml

func main() {
	m := make(yml)

	log.Println("Fetching openapi.yml from typesense api spec")
	err := fetchOpenAPISpec()
	if err != nil {
		log.Fatalf("Aboring: %s", err.Error())
	}

	configFile, err := os.Open("./typesense/api/generator/openapi.yml")
	if err != nil {
		log.Fatalf("Unable to open config file: %s", err.Error())
		return
	}

	err = yaml.NewDecoder(configFile).Decode(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unwrapping the search parameters
	log.Println("Unwrapping search parameters and multi_search parameters")
	unwrapSearchParameters(&m)
	unwrapMultiSearchParameters(&m)
	// Unwrapping import and export parameters
	log.Println("Unwrapping documents import parameters")
	unwrapImportDocuments(&m)
	log.Println("Unwrapping documents export parameters")
	unwrapExportDocuments(&m)
	// Unwrapping update documents with condition parameters
	log.Println("Unwrapping documents update with condition parameters")
	unwrapUpdateDocumentsWithConditionParameters(&m)
	// Unwrapping delete document parameters
	log.Println("Unwrapping documents delete parameters")
	unwrapDeleteDocument(&m)
	log.Println("Unwrapping collections get parameters")
	unwrapGetCollections(&m)
	// Remove additionalProperties from SearchResultHit -> document
	log.Println("Removing additionalProperties from SearchResultHit")
	searchResultHit(&m)

	log.Println("Writing updated spec to generator.yml")
	generatorFile, err := os.Create("./typesense/api/generator/generator.yml")
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

	// Use generator.yml to generate client_gen.go and types_gen.go
	log.Println("Generating client")
	oAPICodeGen()
	log.Println("Successfully Completed !")
}

func fetchOpenAPISpec() error {
	url := "https://raw.githubusercontent.com/typesense/typesense-api-spec/master/openapi.yml"

	// Fetch the spec
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Unable to fetch spec: %s", err.Error())
		return nil
	}
	defer resp.Body.Close()

	// Write the spec to openapi.yml file
	openapiFile, err := os.Create("./typesense/api/generator/openapi.yml")
	if err != nil {
		log.Printf("Unable to write openapi.yml file: %s", err.Error())
		return nil
	}
	defer openapiFile.Close()

	// Write the body to generator.yml file
	_, err = io.Copy(openapiFile, resp.Body)
	return err
}

func searchResultHit(m *yml) {
	properties := (*m)["components"].(yml)["schemas"].(yml)["SearchResultHit"].(yml)["properties"].(yml)
	document := properties["document"].(yml)
	delete(document, "additionalProperties")
}

func unwrapDeleteDocument(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents"].(yml)["delete"].(yml)["parameters"].([]interface{})
	deleteParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)
	for _, obj := range sortedSlice(deleteParameters) {
		newMap := make(yml)
		newMap["name"] = obj.Key
		newMap["in"] = query
		newMap["schema"] = make(yml)
		newMap["schema"].(yml)["type"] = obj.Value.(yml)["type"].(string)
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents"].(yml)["delete"].(yml)["parameters"] = parameters
}

func unwrapUpdateDocumentsWithConditionParameters(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents"].(yml)["patch"].(yml)["parameters"].([]interface{})
	updateParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)
	for _, obj := range sortedSlice(updateParameters) {
		newMap := make(yml)
		newMap["name"] = obj.Key
		newMap["in"] = query
		newMap["schema"] = make(yml)
		newMap["schema"].(yml)["type"] = obj.Value.(yml)["type"].(string)
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents"].(yml)["patch"].(yml)["parameters"] = parameters
}

func unwrapExportDocuments(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents/export"].(yml)["get"].(yml)["parameters"].([]interface{})
	exportParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)
	for _, obj := range sortedSlice(exportParameters) {
		newMap := make(yml)
		newMap["name"] = obj.Key
		newMap["in"] = query
		newMap["schema"] = make(yml)
		if obj.Value.(yml)["type"].(string) == array {
			newMap["schema"].(yml)["type"] = array
			newMap["schema"].(yml)["items"] = obj.Value.(yml)["items"]
		} else {
			newMap["schema"].(yml)["type"] = obj.Value.(yml)["type"].(string)
		}
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents/export"].(yml)["get"].(yml)["parameters"] = parameters
}

func unwrapImportDocuments(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents/import"].(yml)["post"].(yml)["parameters"].([]interface{})
	importParameters := parameters[1].(yml)["schema"].(yml)["properties"].(yml)

	for _, obj := range sortedSlice(importParameters) {
		newMap := make(yml)
		newMap["name"] = obj.Key
		newMap["in"] = query
		newMap["schema"] = make(yml)
		switch {
		// if the param is referencing a schema
		case obj.Value.(yml)["type"] == nil:
			newMap["schema"].(yml)["$ref"] = obj.Value.(yml)["$ref"].(string)
		case obj.Value.(yml)["type"].(string) == array:
			newMap["schema"].(yml)["type"] = array
			newMap["schema"].(yml)["items"] = obj.Value.(yml)["items"]
		default:
			newMap["schema"].(yml)["type"] = obj.Value.(yml)["type"].(string)
		}
		if obj.Value.(yml)["enum"] != nil {
			newMap["schema"].(yml)["enum"] = obj.Value.(yml)["enum"]
		}
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents/import"].(yml)["post"].(yml)["parameters"] = parameters
}

func getSearchParameters(m *yml) yml {
	search := (*m)["components"].(yml)["schemas"].(yml)["SearchParameters"].(yml)["properties"].(yml)
	return search
}

func unwrapSearchParameters(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections/{collectionName}/documents/search"].(yml)["get"].(yml)["parameters"].([]interface{})
	searchParameters := getSearchParameters(m)

	for _, obj := range sortedSlice(searchParameters) {
		newMap := make(yml)
		newMap["name"] = obj.Key
		newMap["in"] = query
		newMap["schema"] = make(yml)
		if obj.Value.(yml)["oneOf"] == nil {
			switch {
			// if the param is referencing a schema
			case obj.Value.(yml)["type"] == nil:
				newMap["schema"].(yml)["$ref"] = obj.Value.(yml)["$ref"].(string)
			case obj.Value.(yml)["type"].(string) == array:
				newMap["schema"].(yml)["type"] = array
				newMap["schema"].(yml)["items"] = obj.Value.(yml)["items"]
			default:
				newMap["schema"].(yml)["type"] = obj.Value.(yml)["type"].(string)
			}
		} else {
			newMap["schema"].(yml)["oneOf"] = obj.Value.(yml)["oneOf"]
		}
		parameters = append(parameters, newMap)
	}

	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(yml)["/collections/{collectionName}/documents/search"].(yml)["get"].(yml)["parameters"] = parameters
}

func unwrapMultiSearchParameters(m *yml) {
	parameters := (*m)["paths"].(yml)["/multi_search"].(yml)["post"].(yml)["parameters"].([]interface{})
	searchParameters := getSearchParameters(m)

	for _, obj := range sortedSlice(searchParameters) {
		newMap := make(yml)
		newMap["name"] = obj.Key
		newMap["in"] = query
		newMap["schema"] = make(yml)
		if obj.Value.(yml)["oneOf"] == nil {
			switch {
			// if the param is referencing a schema
			case obj.Value.(yml)["type"] == nil:
				newMap["schema"].(yml)["$ref"] = obj.Value.(yml)["$ref"].(string)
			case obj.Value.(yml)["type"].(string) == array:
				newMap["schema"].(yml)["type"] = array
				newMap["schema"].(yml)["items"] = obj.Value.(yml)["items"]
			default:
				newMap["schema"].(yml)["type"] = obj.Value.(yml)["type"].(string)
			}
		} else {
			newMap["schema"].(yml)["oneOf"] = obj.Value.(yml)["oneOf"]
		}
		parameters = append(parameters, newMap)
	}

	parameters = parameters[1:]
	(*m)["paths"].(yml)["/multi_search"].(yml)["post"].(yml)["parameters"] = parameters
}

func unwrapGetCollections(m *yml) {
	parameters := (*m)["paths"].(yml)["/collections"].(yml)["get"].(yml)["parameters"].([]interface{})
	deleteParameters := parameters[0].(yml)["schema"].(yml)["properties"].(yml)
	for _, obj := range sortedSlice(deleteParameters) {
		newMap := make(yml)
		newMap["name"] = obj.Key
		newMap["in"] = query
		newMap["schema"] = make(yml)
		newMap["schema"].(yml)["type"] = obj.Value.(yml)["type"].(string)
		// newMap["schema"].(yml)["description"] = obj.Value.(yml)["description"].(string)
		parameters = append(parameters, newMap)
	}
	parameters = parameters[1:]
	(*m)["paths"].(yml)["/collections"].(yml)["get"].(yml)["parameters"] = parameters
}

func oAPICodeGen() {
	cmd := exec.Command("pwd")
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("Unable to get current directory: %s", err.Error())
	}
	currentDir := strings.TrimSpace(string(stdout))
	log.Printf("Current directory: %s", currentDir)
	log.Println("Generating client_gen.go and types_gen.go")

	// Generate client_gen.go and types_gen.go
	err = exec.Command("sh", "./generator.sh").Run()
	if err != nil {
		log.Printf("Error generating client_gen.go and types_gen.go: %s", err.Error())
	}
}
