package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

func traverseMap(m map[string]interface{}, mainKey, prefix string, forTesting bool) {
	for key, value := range m {
		switch v := value.(type) {
		case map[interface{}]interface{}:
			// Convert map[interface{}]interface{} to map[string]interface{}
			nestedMap := make(map[string]interface{})
			for k, vv := range v {
				nestedMap[fmt.Sprintf("%v", k)] = vv
			}
			// Recursive call for nested maps
			newPrefix := fmt.Sprintf("%s.%s", prefix, key)
			traverseMap(nestedMap, mainKey, newPrefix, forTesting)
		case map[string]interface{}:
			// Recursive call for nested maps
			newPrefix := fmt.Sprintf("%s.%s", prefix, key)
			traverseMap(v, mainKey, newPrefix, forTesting)
		default:
			// Generate Tekton result paths for non-map values
			resultPath := fmt.Sprintf("%s_%s.path", prefix, key)
			resultValue := fmt.Sprintf("%v", v)
			if forTesting {
				fmt.Printf("%s = %s\n", resultPath, resultValue)
			} else {
				fmt.Printf("printf \"%s\" \"${%s}\" > \"$(%s)\"\n", resultValue, resultPath, resultPath)
			}
		}
	}
}

func main() {
	// Get the environment variables
	filename := os.Getenv("FILE")
	mainKey := os.Getenv("MAIN_KEY")
	forTesting := os.Getenv("FOR_TESTING")

	// Check if any of the environment variables are empty
	if filename == "" || mainKey == "" || forTesting == "" {
		fmt.Println("Error: Environment variables FILE, MAIN_KEY, and FOR_TESTING must be set")
		os.Exit(1)
	}

	// Convert forTesting to bool
	forTestingBool := false
	if strings.ToLower(forTesting) == "true" {
		forTestingBool = true
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var result map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	mainMap, ok := result[mainKey].(map[interface{}]interface{})
	if !ok {
		log.Fatalf("Main key '%s' not found or not a map", mainKey)
	}

	// Convert map[interface{}]interface{} to map[string]interface{}
	mainMapStr := make(map[string]interface{})
	for k, v := range mainMap {
		mainMapStr[fmt.Sprintf("%v", k)] = v
	}

	traverseMap(mainMapStr, mainKey, "results", forTestingBool)
}
