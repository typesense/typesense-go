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

	// Unwrapping the searchPaths
	unwrapSearchParameters(&m)

	// fmt.Printf("%v", m["paths"].(map[string]interface{})["/collections/{collectionName}/documents/search"].(map[string]interface{})["get"].(map[string]interface{})["parameters"])
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
