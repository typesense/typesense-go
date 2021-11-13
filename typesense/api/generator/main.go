package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	m := make(map[string]interface{})

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

func searchResultHit(m *map[string]interface{}) {
	properties := (*m)["components"].(map[string]interface{})["schemas"].(map[string]interface{})["SearchResultHit"].(map[string]interface{})["properties"].(map[string]interface{})
	document := properties["document"].(map[string]interface{})
	delete(document, "additionalProperties")
}

func unwrapExportDocuments(m *map[string]interface{}) {
	parameters := (*m)["paths"].(map[string]interface{})["/collections/{collectionName}/documents/export"].(map[string]interface{})["get"].(map[string]interface{})["parameters"].([]interface{})
	exportParameters := parameters[1].(map[string]interface{})["schema"].(map[string]interface{})["properties"].(map[string]interface{})
	for k, v := range exportParameters {
		newMap := make(map[string]interface{})
		newMap["name"] = k
		newMap["in"] = "query"
		newMap["schema"] = make(map[string]interface{})
		if v.(map[string]interface{})["type"].(string) == "array" {
			newMap["schema"].(map[string]interface{})["type"] = "array"
			newMap["schema"].(map[string]interface{})["items"] = v.(map[string]interface{})["items"]
		} else {
			newMap["schema"].(map[string]interface{})["type"] = v.(map[string]interface{})["type"].(string)
		}
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(map[string]interface{})["/collections/{collectionName}/documents/export"].(map[string]interface{})["get"].(map[string]interface{})["parameters"] = parameters
}

func unwrapImportDocuments(m *map[string]interface{}) {
	parameters := (*m)["paths"].(map[string]interface{})["/collections/{collectionName}/documents/import"].(map[string]interface{})["post"].(map[string]interface{})["parameters"].([]interface{})
	importParameters := parameters[1].(map[string]interface{})["schema"].(map[string]interface{})["properties"].(map[string]interface{})

	for k, v := range importParameters {
		newMap := make(map[string]interface{})
		newMap["name"] = k
		newMap["in"] = "query"
		newMap["schema"] = make(map[string]interface{})
		newMap["schema"].(map[string]interface{})["type"] = v.(map[string]interface{})["type"].(string)
		if v.(map[string]interface{})["enum"] != nil {
			newMap["schema"].(map[string]interface{})["enum"] = v.(map[string]interface{})["enum"]
		}
		parameters = append(parameters, newMap)
	}
	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(map[string]interface{})["/collections/{collectionName}/documents/import"].(map[string]interface{})["post"].(map[string]interface{})["parameters"] = parameters
}

func unwrapSearchParameters(m *map[string]interface{}) {
	parameters := (*m)["paths"].(map[string]interface{})["/collections/{collectionName}/documents/search"].(map[string]interface{})["get"].(map[string]interface{})["parameters"].([]interface{})
	searchParameters := parameters[1].(map[string]interface{})["schema"].(map[string]interface{})["properties"].(map[string]interface{})

	for k, v := range searchParameters {
		newMap := make(map[string]interface{})
		newMap["name"] = k
		if k == "q" || k == "query_by" {
			newMap["required"] = true
		}
		newMap["in"] = "query"
		newMap["schema"] = make(map[string]interface{})
		if v.(map[string]interface{})["oneOf"] == nil {
			if v.(map[string]interface{})["type"].(string) == "array" {
				newMap["schema"].(map[string]interface{})["type"] = "array"
				newMap["schema"].(map[string]interface{})["items"] = v.(map[string]interface{})["items"]
			} else {
				newMap["schema"].(map[string]interface{})["type"] = v.(map[string]interface{})["type"].(string)
			}
		} else {
			newMap["schema"].(map[string]interface{})["oneOf"] = v.(map[string]interface{})["oneOf"]
		}
		parameters = append(parameters, newMap)
	}

	parameters = append(parameters[:1], parameters[2:]...)
	(*m)["paths"].(map[string]interface{})["/collections/{collectionName}/documents/search"].(map[string]interface{})["get"].(map[string]interface{})["parameters"] = parameters
}
